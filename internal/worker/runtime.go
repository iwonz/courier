package worker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

const serverReadTimeout = 5 * time.Second

type leaseEntry struct {
	deliveryID delivery.ID
	connection net.Conn
}

type Runtime struct {
	store         *delivery.Store
	server        delivery.Server
	compatibility string
	now           func() time.Time

	mutex        sync.Mutex
	deliveries   map[delivery.ID]delivery.Delivery
	ownedTemps   map[delivery.ID][]delivery.ID
	leases       map[delivery.ID]leaseEntry
	keepalives   map[delivery.ID]bool
	subscribers  map[delivery.ID]chan ProgressEvent
	hadDelivery  bool
	listener     net.Listener
	shutdown     chan struct{}
	shutdownOnce sync.Once
	connections  sync.WaitGroup
	host         DeliveryHost
}

func NewRuntime(ctx context.Context, store *delivery.Store, server delivery.Server, compatibility string) (*Runtime, error) {
	if store == nil || server.State != delivery.StateStarting || invalidCompatibility(compatibility) {
		return nil, fmt.Errorf("%w: invalid worker runtime", delivery.ErrInvalid)
	}
	runtime := &Runtime{
		store: store, server: server, compatibility: compatibility, now: time.Now,
		deliveries: map[delivery.ID]delivery.Delivery{}, ownedTemps: map[delivery.ID][]delivery.ID{}, leases: map[delivery.ID]leaseEntry{},
		keepalives: map[delivery.ID]bool{}, subscribers: map[delivery.ID]chan ProgressEvent{}, shutdown: make(chan struct{}),
	}
	if _, err := updateStore(ctx, store, func(registry *delivery.Registry) error {
		if err := registry.RegisterServer(server); err != nil {
			return err
		}
		return registry.TransitionServer(server.ID, delivery.StateActive, server.UpdatedAt)
	}); err != nil {
		return nil, err
	}
	runtime.server.State = delivery.StateActive
	return runtime, nil
}

func (runtime *Runtime) ServerID() delivery.ID { return runtime.server.ID }

func (runtime *Runtime) Done() <-chan struct{} { return runtime.shutdown }

func (runtime *Runtime) HasDeliveries() bool {
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	return len(runtime.deliveries) != 0
}

func (runtime *Runtime) AttachHost(host DeliveryHost) error {
	if host == nil {
		return fmt.Errorf("%w: delivery host is required", delivery.ErrInvalid)
	}
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	if runtime.host != nil || runtime.listener != nil || len(runtime.deliveries) != 0 {
		return fmt.Errorf("%w: delivery host cannot be attached", delivery.ErrInvalid)
	}
	runtime.host = host
	return nil
}

func (runtime *Runtime) Run(listener net.Listener) error {
	if listener == nil {
		return fmt.Errorf("%w: control listener is required", delivery.ErrInvalid)
	}
	runtime.mutex.Lock()
	if runtime.listener != nil {
		runtime.mutex.Unlock()
		return fmt.Errorf("%w: worker runtime is already running", delivery.ErrInvalid)
	}
	runtime.listener = listener
	if runtime.isShuttingDown() {
		runtime.mutex.Unlock()
		_ = listener.Close()
		return nil
	}
	runtime.mutex.Unlock()

	for {
		connection, err := listener.Accept()
		if err != nil {
			select {
			case <-runtime.shutdown:
				runtime.connections.Wait()
				return nil
			default:
				stopErr := runtime.StopServer(context.Background(), delivery.ReasonFailed)
				runtime.connections.Wait()
				return errors.Join(err, stopErr)
			}
		}
		runtime.connections.Add(1)
		go func() {
			defer runtime.connections.Done()
			runtime.handle(connection)
		}()
	}
}

func (runtime *Runtime) Register(ctx context.Context, request RegisterRequest, connection net.Conn) (RegisterResponse, error) {
	if err := request.Validate(); err != nil {
		return RegisterResponse{}, err
	}
	if request.Delivery.ServerID != runtime.server.ID || request.Delivery.State != delivery.StateStarting {
		return RegisterResponse{}, fmt.Errorf("%w: delivery does not belong to this starting worker", delivery.ErrInvalid)
	}
	if request.LeaseID != "" && connection == nil {
		return RegisterResponse{}, fmt.Errorf("%w: foreground registration requires a connection", delivery.ErrInvalid)
	}
	if len(request.RuntimeDefinition) != 0 && runtime.host == nil {
		return RegisterResponse{}, fmt.Errorf("%w: runtime definition requires a delivery host", delivery.ErrInvalid)
	}

	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	if runtime.isShuttingDown() {
		return RegisterResponse{}, fmt.Errorf("%w: worker is stopping", delivery.ErrInvalid)
	}
	if _, exists := runtime.deliveries[request.Delivery.ID]; exists {
		return RegisterResponse{}, fmt.Errorf("%w: delivery %s", delivery.ErrDuplicate, request.Delivery.ID)
	}
	if request.LeaseID != "" {
		if _, exists := runtime.leases[request.LeaseID]; exists {
			return RegisterResponse{}, fmt.Errorf("%w: lease %s", delivery.ErrDuplicate, request.LeaseID)
		}
	}
	item := request.Delivery
	hostRegistered := false
	var ownedTemps []delivery.OwnedTemp
	if runtime.host != nil && len(request.RuntimeDefinition) != 0 {
		if err := runtime.host.Register(ctx, item, request.RuntimeDefinition); err != nil {
			return RegisterResponse{}, err
		}
		hostRegistered = true
		ownedTemps = runtime.host.OwnedTemps(item.ID)
	}
	if _, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		if err := registry.RegisterDelivery(item); err != nil {
			return err
		}
		for _, temporary := range ownedTemps {
			if err := registry.AddOwnedTemp(temporary); err != nil {
				return err
			}
		}
		return registry.TransitionDelivery(item.ID, delivery.StateActive, item.UpdatedAt)
	}); err != nil {
		if hostRegistered {
			err = errors.Join(err, runtime.host.Stop(context.Background(), item.ID))
		}
		return RegisterResponse{}, err
	}
	item.State = delivery.StateActive
	runtime.deliveries[item.ID] = item
	for _, temporary := range ownedTemps {
		runtime.ownedTemps[item.ID] = append(runtime.ownedTemps[item.ID], temporary.ID)
	}
	if request.LeaseID != "" {
		runtime.leases[request.LeaseID] = leaseEntry{deliveryID: item.ID, connection: connection}
	}
	runtime.hadDelivery = true
	runtime.publishLocked()
	return RegisterResponse{ServerID: runtime.server.ID, DeliveryID: item.ID, LeaseID: request.LeaseID}, nil
}

func (runtime *Runtime) ClaimLease(request LeaseRequest, connection net.Conn) error {
	if err := request.Validate(); err != nil {
		return err
	}
	if connection == nil {
		return fmt.Errorf("%w: lease connection is required", delivery.ErrInvalid)
	}
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	if _, exists := runtime.deliveries[request.DeliveryID]; !exists {
		return fmt.Errorf("%w: delivery %s", delivery.ErrNotFound, request.DeliveryID)
	}
	if runtime.isShuttingDown() {
		return fmt.Errorf("%w: worker is stopping", delivery.ErrInvalid)
	}
	if _, exists := runtime.leases[request.LeaseID]; exists || runtime.deliveryHasLeaseLocked(request.DeliveryID) {
		return fmt.Errorf("%w: delivery or lease is already owned", delivery.ErrDuplicate)
	}
	runtime.leases[request.LeaseID] = leaseEntry{deliveryID: request.DeliveryID, connection: connection}
	return nil
}

func (runtime *Runtime) ReleaseLease(ctx context.Context, request ReleaseLeaseRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	runtime.mutex.Lock()
	entry, exists := runtime.leases[request.LeaseID]
	if !exists || entry.deliveryID != request.DeliveryID {
		runtime.mutex.Unlock()
		return fmt.Errorf("%w: lease %s", delivery.ErrNotFound, request.LeaseID)
	}
	delete(runtime.leases, request.LeaseID)
	runtime.mutex.Unlock()
	_, err := runtime.StopDelivery(ctx, request.DeliveryID, delivery.ReasonStopped)
	return err
}

func (runtime *Runtime) UpdatePolicy(ctx context.Context, request UpdatePolicyRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	item, exists := runtime.deliveries[request.DeliveryID]
	if !exists {
		return fmt.Errorf("%w: delivery %s", delivery.ErrNotFound, request.DeliveryID)
	}
	var commitPolicy, cancelPolicy func()
	if runtime.host != nil {
		var err error
		commitPolicy, cancelPolicy, err = runtime.host.PreparePolicy(ctx, request.DeliveryID, request.ExpectedVersion, request.Policy)
		if err != nil {
			return err
		}
		defer cancelPolicy()
	}
	at := runtime.now().UTC()
	if _, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		return registry.UpdatePolicy(request.DeliveryID, request.ExpectedVersion, request.Policy, at)
	}); err != nil {
		return err
	}
	if commitPolicy != nil {
		commitPolicy()
	}
	item.Policy = request.Policy
	item.UpdatedAt = at
	runtime.deliveries[item.ID] = item
	runtime.publishLocked()
	return nil
}

func (runtime *Runtime) SetCounters(ctx context.Context, id delivery.ID, counters delivery.CounterSnapshot) error {
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	item, exists := runtime.deliveries[id]
	if !exists {
		return fmt.Errorf("%w: delivery %s", delivery.ErrNotFound, id)
	}
	at := runtime.now().UTC()
	if _, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		return registry.SetCounters(id, counters, at)
	}); err != nil {
		return err
	}
	item.Counters = counters
	item.UpdatedAt = at
	runtime.deliveries[id] = item
	runtime.publishLocked()
	return nil
}

func (runtime *Runtime) StopDelivery(ctx context.Context, id delivery.ID, reason delivery.TombstoneReason) (bool, error) {
	runtime.mutex.Lock()
	item, exists := runtime.deliveries[id]
	if !exists {
		runtime.mutex.Unlock()
		return false, nil
	}
	at := runtime.now().UTC()
	tombstone := delivery.Tombstone{ID: delivery.NewID(), TargetID: id, Kind: delivery.TargetDelivery, Reason: reason, At: at}
	last := len(runtime.deliveries) == 1 && len(runtime.keepalives) == 0
	serverTombstone := delivery.Tombstone{}
	if last {
		serverTombstone = delivery.Tombstone{ID: delivery.NewID(), TargetID: runtime.server.ID, Kind: delivery.TargetServer, Reason: reason, At: at}
	}
	_, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		if err := finishDelivery(registry, item, at); err != nil {
			return err
		}
		if err := registry.TombstoneDelivery(id, tombstone); err != nil {
			return err
		}
		if last {
			if err := finishServer(registry, runtime.server, at); err != nil {
				return err
			}
			return registry.TombstoneServer(runtime.server.ID, serverTombstone)
		}
		return nil
	})
	if err != nil {
		runtime.mutex.Unlock()
		return false, err
	}
	temporaryIDs := append([]delivery.ID(nil), runtime.ownedTemps[id]...)
	var cleanupErr error
	if runtime.host != nil {
		cleanupErr = runtime.host.Stop(ctx, id)
	}
	if cleanupErr == nil && len(temporaryIDs) != 0 {
		cleanupErr = runtime.removeOwnedTemps(ctx, temporaryIDs)
	}
	connections := runtime.removeDeliveryLocked(id)
	if last {
		runtime.server.State = delivery.StateStopped
		runtime.requestShutdownLocked()
	}
	runtime.publishLocked()
	runtime.mutex.Unlock()
	closeConnections(connections)
	return true, cleanupErr
}

func (runtime *Runtime) StopServer(ctx context.Context, reason delivery.TombstoneReason) error {
	runtime.mutex.Lock()
	if runtime.server.State == delivery.StateStopped {
		runtime.mutex.Unlock()
		return nil
	}
	at := runtime.now().UTC()
	items := make([]delivery.Delivery, 0, len(runtime.deliveries))
	for _, item := range runtime.deliveries {
		items = append(items, item)
	}
	_, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		for _, item := range items {
			if err := finishDelivery(registry, item, at); err != nil {
				return err
			}
			if err := registry.TombstoneDelivery(item.ID, delivery.Tombstone{ID: delivery.NewID(), TargetID: item.ID, Kind: delivery.TargetDelivery, Reason: reason, At: at}); err != nil {
				return err
			}
		}
		if err := finishServer(registry, runtime.server, at); err != nil {
			return err
		}
		return registry.TombstoneServer(runtime.server.ID, delivery.Tombstone{ID: delivery.NewID(), TargetID: runtime.server.ID, Kind: delivery.TargetServer, Reason: reason, At: at})
	})
	if err != nil {
		runtime.mutex.Unlock()
		return err
	}
	var cleanupErr error
	cleanedTemps := make([]delivery.ID, 0)
	for _, item := range items {
		var stopErr error
		if runtime.host != nil {
			stopErr = runtime.host.Stop(ctx, item.ID)
		}
		cleanupErr = errors.Join(cleanupErr, stopErr)
		if stopErr == nil {
			cleanedTemps = append(cleanedTemps, runtime.ownedTemps[item.ID]...)
		}
	}
	if len(cleanedTemps) != 0 {
		tempErr := runtime.removeOwnedTemps(ctx, cleanedTemps)
		cleanupErr = errors.Join(cleanupErr, tempErr)
	}
	connections := make([]net.Conn, 0, len(runtime.leases))
	for _, lease := range runtime.leases {
		connections = append(connections, lease.connection)
	}
	runtime.deliveries = map[delivery.ID]delivery.Delivery{}
	runtime.ownedTemps = map[delivery.ID][]delivery.ID{}
	runtime.leases = map[delivery.ID]leaseEntry{}
	runtime.keepalives = map[delivery.ID]bool{}
	runtime.server.State = delivery.StateStopped
	runtime.requestShutdownLocked()
	runtime.publishLocked()
	runtime.mutex.Unlock()
	closeConnections(connections)
	return cleanupErr
}

func (runtime *Runtime) AcquireKeepalive() (delivery.ID, error) {
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	if runtime.isShuttingDown() {
		return "", fmt.Errorf("%w: worker is stopping", delivery.ErrInvalid)
	}
	id := delivery.NewID()
	runtime.keepalives[id] = true
	return id, nil
}

func (runtime *Runtime) ReleaseKeepalive(ctx context.Context, id delivery.ID) error {
	runtime.mutex.Lock()
	if !runtime.keepalives[id] {
		runtime.mutex.Unlock()
		return fmt.Errorf("%w: keepalive %s", delivery.ErrNotFound, id)
	}
	delete(runtime.keepalives, id)
	shouldStop := runtime.hadDelivery && len(runtime.deliveries) == 0 && len(runtime.keepalives) == 0
	runtime.mutex.Unlock()
	if shouldStop {
		return runtime.StopServer(ctx, delivery.ReasonCompleted)
	}
	return nil
}

func (runtime *Runtime) Snapshot(ctx context.Context) (delivery.Snapshot, error) {
	return runtime.store.Load(ctx)
}

func (runtime *Runtime) Subscribe(id delivery.ID) (<-chan ProgressEvent, func(), error) {
	runtime.mutex.Lock()
	if id != "" {
		if _, exists := runtime.deliveries[id]; !exists {
			runtime.mutex.Unlock()
			return nil, nil, fmt.Errorf("%w: delivery %s", delivery.ErrNotFound, id)
		}
	}
	if runtime.isShuttingDown() {
		runtime.mutex.Unlock()
		return nil, nil, fmt.Errorf("%w: worker is stopping", delivery.ErrInvalid)
	}
	subscriptionID := delivery.NewID()
	channel := make(chan ProgressEvent, 1)
	runtime.subscribers[subscriptionID] = channel
	runtime.publishOneLocked(channel, id)
	runtime.mutex.Unlock()
	var once sync.Once
	cancel := func() {
		once.Do(func() {
			runtime.mutex.Lock()
			if _, exists := runtime.subscribers[subscriptionID]; exists {
				delete(runtime.subscribers, subscriptionID)
				close(channel)
			}
			runtime.mutex.Unlock()
		})
	}
	return channel, cancel, nil
}

func (runtime *Runtime) handle(connection net.Conn) {
	defer connection.Close()
	_ = connection.SetReadDeadline(time.Now().Add(serverReadTimeout))
	var request ipc.Request
	if err := ipc.ReadFrame(connection, &request); err != nil || request.Validate() != nil {
		return
	}
	switch request.Operation {
	case ipc.OperationHello:
		runtime.handleHello(connection, request)
	case ipc.OperationRegister:
		runtime.handleRegister(connection, request)
	case ipc.OperationLease:
		runtime.handleLease(connection, request)
	case ipc.OperationList:
		snapshot, err := runtime.Snapshot(context.Background())
		runtime.respond(connection, request, ListResponse{Snapshot: snapshot}, err)
	case ipc.OperationUpdatePolicy:
		var payload UpdatePolicyRequest
		err := ipc.DecodePayload(request.Payload, &payload)
		if err == nil {
			err = runtime.UpdatePolicy(context.Background(), payload)
		}
		runtime.respond(connection, request, nil, err)
	case ipc.OperationStopDelivery:
		var payload TargetRequest
		err := ipc.DecodePayload(request.Payload, &payload)
		if err == nil {
			err = payload.Validate()
		}
		if err == nil {
			_, err = runtime.StopDelivery(context.Background(), payload.ID, delivery.ReasonStopped)
		}
		runtime.respond(connection, request, nil, err)
	case ipc.OperationStopServer, ipc.OperationShutdown:
		err := runtime.StopServer(context.Background(), delivery.ReasonStopped)
		runtime.respond(connection, request, nil, err)
	case ipc.OperationSubscribeProgress:
		runtime.handleSubscription(connection, request)
	case ipc.OperationReleaseLease:
		runtime.respond(connection, request, nil, fmt.Errorf("%w: release requires a lease connection", delivery.ErrInvalid))
	}
}

func (runtime *Runtime) handleHello(connection net.Conn, request ipc.Request) {
	var payload HelloRequest
	err := ipc.DecodePayload(request.Payload, &payload)
	if err == nil {
		err = payload.Validate()
	}
	if err == nil && payload.ServerID != runtime.server.ID {
		err = fmt.Errorf("%w: server identity mismatch", delivery.ErrNotFound)
	}
	if err == nil && payload.Compatibility != runtime.compatibility {
		err = fmt.Errorf("%w: worker compatibility mismatch", delivery.ErrRevisionConflict)
	}
	runtime.respond(connection, request, HelloResponse{ServerID: runtime.server.ID, Compatibility: runtime.compatibility}, err)
}

func (runtime *Runtime) handleRegister(connection net.Conn, request ipc.Request) {
	var payload RegisterRequest
	err := ipc.DecodePayload(request.Payload, &payload)
	if err == nil {
		err = payload.Validate()
	}
	var result RegisterResponse
	if err == nil {
		result, err = runtime.Register(context.Background(), payload, connection)
	}
	responded := runtime.respond(connection, request, result, err)
	if !responded && err == nil && payload.LeaseID != "" {
		_ = runtime.releaseLostLease(context.Background(), payload.Delivery.ID, payload.LeaseID)
	}
	if !responded || err != nil || payload.LeaseID == "" {
		return
	}
	runtime.holdLease(connection, payload.Delivery.ID, payload.LeaseID)
}

func (runtime *Runtime) handleLease(connection net.Conn, request ipc.Request) {
	var payload LeaseRequest
	err := ipc.DecodePayload(request.Payload, &payload)
	if err == nil {
		err = payload.Validate()
	}
	if err == nil {
		err = runtime.ClaimLease(payload, connection)
	}
	if !runtime.respond(connection, request, nil, err) || err != nil {
		return
	}
	runtime.holdLease(connection, payload.DeliveryID, payload.LeaseID)
}

func (runtime *Runtime) holdLease(connection net.Conn, deliveryID, leaseID delivery.ID) {
	_ = connection.SetReadDeadline(time.Time{})
	var request ipc.Request
	if err := ipc.ReadFrame(connection, &request); err == nil && request.Validate() == nil && request.Operation == ipc.OperationReleaseLease {
		var payload ReleaseLeaseRequest
		err = ipc.DecodePayload(request.Payload, &payload)
		if err == nil && (payload.DeliveryID != deliveryID || payload.LeaseID != leaseID) {
			err = fmt.Errorf("%w: lease release mismatch", delivery.ErrInvalid)
		}
		if err == nil {
			err = runtime.ReleaseLease(context.Background(), payload)
		}
		runtime.respond(connection, request, nil, err)
		return
	}
	_ = runtime.releaseLostLease(context.Background(), deliveryID, leaseID)
}

func (runtime *Runtime) handleSubscription(connection net.Conn, request ipc.Request) {
	var payload ProgressRequest
	err := ipc.DecodePayload(request.Payload, &payload)
	if err == nil {
		err = payload.Validate()
	}
	var events <-chan ProgressEvent
	var cancel func()
	if err == nil {
		events, cancel, err = runtime.Subscribe(payload.DeliveryID)
	}
	if err != nil {
		runtime.respond(connection, request, nil, err)
		return
	}
	defer cancel()
	_ = connection.SetDeadline(time.Time{})
	peerClosed := make(chan struct{})
	go func() {
		buffer := []byte{0}
		_, _ = connection.Read(buffer)
		close(peerClosed)
	}()
	for {
		select {
		case event, open := <-events:
			if !open {
				return
			}
			response, responseErr := ipc.Success(request, event)
			if responseErr != nil || ipc.WriteFrame(connection, response) != nil {
				return
			}
		case <-peerClosed:
			return
		}
	}
}

func (runtime *Runtime) respond(connection net.Conn, request ipc.Request, payload any, err error) bool {
	var response ipc.Response
	var buildErr error
	if err == nil {
		response, buildErr = ipc.Success(request, payload)
	} else {
		response, buildErr = ipc.Rejection(request, errorCode(err), publicError(err))
	}
	return buildErr == nil && ipc.WriteFrame(connection, response) == nil
}

func (runtime *Runtime) releaseLostLease(ctx context.Context, deliveryID, leaseID delivery.ID) error {
	runtime.mutex.Lock()
	entry, exists := runtime.leases[leaseID]
	if !exists || entry.deliveryID != deliveryID {
		runtime.mutex.Unlock()
		return nil
	}
	delete(runtime.leases, leaseID)
	runtime.mutex.Unlock()
	_, err := runtime.StopDelivery(ctx, deliveryID, delivery.ReasonStopped)
	return err
}

func (runtime *Runtime) deliveryHasLeaseLocked(id delivery.ID) bool {
	for _, lease := range runtime.leases {
		if lease.deliveryID == id {
			return true
		}
	}
	return false
}

func (runtime *Runtime) removeDeliveryLocked(id delivery.ID) []net.Conn {
	delete(runtime.deliveries, id)
	delete(runtime.ownedTemps, id)
	connections := []net.Conn{}
	for leaseID, lease := range runtime.leases {
		if lease.deliveryID == id {
			connections = append(connections, lease.connection)
			delete(runtime.leases, leaseID)
		}
	}
	return connections
}

func (runtime *Runtime) removeOwnedTemps(ctx context.Context, ids []delivery.ID) error {
	_, err := updateStore(ctx, runtime.store, func(registry *delivery.Registry) error {
		for _, id := range ids {
			if err := registry.RemoveOwnedTemp(id); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (runtime *Runtime) requestShutdownLocked() {
	runtime.shutdownOnce.Do(func() {
		close(runtime.shutdown)
		if runtime.listener != nil {
			_ = runtime.listener.Close()
		}
		for _, subscriber := range runtime.subscribers {
			close(subscriber)
		}
		runtime.subscribers = map[delivery.ID]chan ProgressEvent{}
	})
}

func (runtime *Runtime) isShuttingDown() bool {
	select {
	case <-runtime.shutdown:
		return true
	default:
		return false
	}
}

func (runtime *Runtime) publishLocked() {
	for _, subscriber := range runtime.subscribers {
		runtime.publishOneLocked(subscriber, "")
	}
}

func (runtime *Runtime) publishOneLocked(channel chan ProgressEvent, filter delivery.ID) {
	snapshot, err := runtime.store.Load(context.Background())
	if err != nil {
		return
	}
	if filter != "" {
		filtered := snapshot.Deliveries[:0]
		for _, item := range snapshot.Deliveries {
			if item.ID == filter {
				filtered = append(filtered, item)
			}
		}
		snapshot.Deliveries = filtered
	}
	event := ProgressEvent{Snapshot: snapshot}
	select {
	case <-channel:
	default:
	}
	select {
	case channel <- event:
	default:
		// Never let a slow or detached subscriber block lifecycle state.
	}
}

func updateStore(ctx context.Context, store *delivery.Store, mutate func(*delivery.Registry) error) (delivery.Snapshot, error) {
	for {
		current, err := store.Load(ctx)
		if err != nil {
			return delivery.Snapshot{}, err
		}
		next, err := store.Update(ctx, current.Revision, mutate)
		if !errors.Is(err, delivery.ErrRevisionConflict) {
			return next, err
		}
		if err := ctx.Err(); err != nil {
			return delivery.Snapshot{}, err
		}
	}
}

func finishDelivery(registry *delivery.Registry, item delivery.Delivery, at time.Time) error {
	switch item.State {
	case delivery.StateStarting, delivery.StateActive:
		if err := registry.TransitionDelivery(item.ID, delivery.StateStopping, at); err != nil {
			return err
		}
		return registry.TransitionDelivery(item.ID, delivery.StateStopped, at)
	case delivery.StateStopping, delivery.StateFailed:
		return registry.TransitionDelivery(item.ID, delivery.StateStopped, at)
	case delivery.StateStopped:
		return nil
	default:
		return fmt.Errorf("%w: invalid delivery terminal state", delivery.ErrInvalid)
	}
}

func finishServer(registry *delivery.Registry, server delivery.Server, at time.Time) error {
	switch server.State {
	case delivery.StateStarting, delivery.StateActive:
		if err := registry.TransitionServer(server.ID, delivery.StateStopping, at); err != nil {
			return err
		}
		return registry.TransitionServer(server.ID, delivery.StateStopped, at)
	case delivery.StateStopping, delivery.StateFailed:
		return registry.TransitionServer(server.ID, delivery.StateStopped, at)
	case delivery.StateStopped:
		return nil
	default:
		return fmt.Errorf("%w: invalid server terminal state", delivery.ErrInvalid)
	}
}

func errorCode(err error) ipc.ErrorCode {
	switch {
	case errors.Is(err, delivery.ErrNotFound):
		return ipc.CodeNotFound
	case errors.Is(err, delivery.ErrDuplicate), errors.Is(err, delivery.ErrRevisionConflict):
		return ipc.CodeConflict
	case errors.Is(err, ipc.ErrProtocol), errors.Is(err, delivery.ErrInvalid):
		return ipc.CodeInvalid
	default:
		return ipc.CodeInternal
	}
}

func publicError(err error) string {
	if err == nil {
		return "operation failed"
	}
	if errorCode(err) == ipc.CodeInternal {
		return "internal worker error"
	}
	message := err.Error()
	if len(message) > 1024 {
		message = message[:1024]
	}
	return message
}

func closeConnections(connections []net.Conn) {
	for _, connection := range connections {
		if connection != nil {
			_ = connection.Close()
		}
	}
}
