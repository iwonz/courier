package delivery

import (
	"errors"
	"math"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/progress"
)

const (
	serverID   ID = "00000000-0000-4000-8000-000000000001"
	deliveryID ID = "00000000-0000-4000-8000-000000000002"
	thirdID    ID = "00000000-0000-4000-8000-000000000003"
	fourthID   ID = "00000000-0000-4000-8000-000000000004"
)

var testTime = time.Unix(1_700_000_000, 0).UTC()

func validServer() Server {
	return Server{ID: serverID, Bind: "127.0.0.1:8080", ControlEndpoint: "control.sock", ProcessID: 42, State: StateStarting, StartedAt: testTime, UpdatedAt: testTime}
}

func validDelivery() Delivery {
	return Delivery{ID: deliveryID, ServerID: serverID, Route: RoutePathToWeb, State: StateStarting, Policy: DefaultPolicy(), CreatedAt: testTime, UpdatedAt: testTime}
}

func TestIdentifiersStatesAndRoutes(t *testing.T) {
	generated := NewID()
	if !generated.Valid() {
		t.Fatalf("generated ID=%q", generated)
	}
	if parsed, err := ParseID(string(serverID)); err != nil || parsed != serverID {
		t.Fatalf("parsed=%q err=%v", parsed, err)
	}
	for _, value := range []string{"bad", "00000000-0000-0000-0000-000000000000", "00000000-0000-4000-8000-000000000001"[:35], "00000000-0000-4000-8000-00000000000A"} {
		if _, err := ParseID(value); !errors.Is(err, ErrInvalid) {
			t.Fatalf("value=%q err=%v", value, err)
		}
	}
	if ID("bad").Valid() {
		t.Fatal("invalid ID accepted")
	}

	states := []State{StateStarting, StateActive, StateStopping, StateStopped, StateFailed}
	for _, state := range states {
		if !state.Valid() {
			t.Fatalf("invalid state %q", state)
		}
	}
	if State("unknown").Valid() {
		t.Fatal("unknown state accepted")
	}
	allowed := map[[2]State]bool{
		{StateStarting, StateActive}: true, {StateStarting, StateStopping}: true, {StateStarting, StateFailed}: true,
		{StateActive, StateStopping}: true, {StateActive, StateFailed}: true,
		{StateStopping, StateStopped}: true, {StateStopping, StateFailed}: true,
		{StateFailed, StateStopped}: true,
	}
	for _, from := range append(states, State("unknown")) {
		for _, to := range states {
			if CanTransition(from, to) != allowed[[2]State{from, to}] {
				t.Fatalf("transition %s to %s", from, to)
			}
		}
	}

	for _, route := range []Route{RouteWebToPath, RoutePathToWeb, RouteWebhookToPath, RoutePathToHTTP} {
		if !route.Valid() {
			t.Fatalf("route=%s", route)
		}
	}
	if Route("path-to-path").Valid() {
		t.Fatal("non-delivery route accepted")
	}
}

func TestPolicyValidation(t *testing.T) {
	policy := DefaultPolicy()
	if err := policy.Validate(); err != nil {
		t.Fatal(err)
	}
	if policy.Version != 1 || policy.AuthAttempts != 5 || !policy.DeliveryLimit.Unlimited || policy.MaxFileSize.Value != 10<<30 || policy.MaxExtractedSize.Value != 100<<30 || !policy.UploadRate.Unlimited || !policy.DownloadRate.Unlimited {
		t.Fatalf("default=%+v", policy)
	}
	for _, auth := range []AuthMode{AuthNone, AuthBasic, AuthPassword} {
		candidate := policy
		candidate.Auth = auth
		if err := candidate.Validate(); err != nil {
			t.Fatalf("auth=%s err=%v", auth, err)
		}
	}
	for _, action := range []AuthFailAction{AuthFailBan, AuthFailStop} {
		candidate := policy
		candidate.AuthFailAction = action
		if err := candidate.Validate(); err != nil {
			t.Fatalf("action=%s err=%v", action, err)
		}
	}
	policy.AllowIP = []string{"127.0.0.1/32", "2001:db8::/32"}
	if err := policy.Validate(); err != nil {
		t.Fatal(err)
	}

	invalid := []func(*Policy){
		func(value *Policy) { value.Version = 0 },
		func(value *Policy) { value.AuthAttempts = 0 },
		func(value *Policy) { value.Auth = "token" },
		func(value *Policy) { value.AuthFailAction = "ignore" },
		func(value *Policy) { value.DeliveryLimit.Value = -1 },
		func(value *Policy) { value.MaxFileSize = Limit{Unlimited: true, Value: 1} },
		func(value *Policy) { value.AllowIP = []string{"bad"} },
		func(value *Policy) { value.AllowIP = []string{"127.0.0.0/24", "127.0.0.0/24"} },
		func(value *Policy) { value.AllowIP = []string{"127.0.0.1/24"} },
	}
	for index, mutate := range invalid {
		candidate := DefaultPolicy()
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("case=%d err=%v", index, err)
		}
	}
}

func TestRecordValidation(t *testing.T) {
	if err := (CounterSnapshot{}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (CounterSnapshot{Read: -1}).Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("counter=%v", err)
	}
	if err := (CounterSnapshot{Read: 1, Sent: 2}).Validate(); !errors.Is(err, ErrInvalid) {
		t.Fatalf("counter ordering=%v", err)
	}
	server := validServer()
	if err := server.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*Server){
		func(value *Server) { value.ID = "bad" }, func(value *Server) { value.Bind = " " }, func(value *Server) { value.Bind = "localhost" },
		func(value *Server) { value.Bind = "localhost:http" }, func(value *Server) { value.Bind = "localhost:0" }, func(value *Server) { value.ControlEndpoint = "bad\nvalue" },
		func(value *Server) { value.ProcessID = 0 }, func(value *Server) { value.State = "bad" }, func(value *Server) { value.StartedAt = time.Time{} },
		func(value *Server) { value.Compatibility = "bad\nvalue" },
	} {
		candidate := server
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("server=%+v err=%v", candidate, err)
		}
	}
	delivery := validDelivery()
	if err := delivery.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*Delivery){
		func(value *Delivery) { value.ID = "bad" }, func(value *Delivery) { value.ServerID = "bad" }, func(value *Delivery) { value.Route = "bad" },
		func(value *Delivery) { value.State = "bad" }, func(value *Delivery) { value.UpdatedAt = value.CreatedAt.Add(-time.Second) },
		func(value *Delivery) { value.Policy.Version = 0 }, func(value *Delivery) { value.Counters.Read = -1 },
		func(value *Delivery) { value.Source = "bad\x00value" },
	} {
		candidate := delivery
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("delivery=%+v err=%v", candidate, err)
		}
	}

	for _, kind := range []TargetKind{TargetServer, TargetDelivery} {
		for _, reason := range []TombstoneReason{ReasonStopped, ReasonCompleted, ReasonFailed, ReasonStale} {
			value := Tombstone{ID: thirdID, TargetID: deliveryID, Kind: kind, Reason: reason, At: testTime}
			if err := value.Validate(); err != nil {
				t.Fatalf("tombstone=%+v err=%v", value, err)
			}
		}
	}
	for _, value := range []Tombstone{
		{ID: "bad", TargetID: deliveryID, Kind: TargetDelivery, Reason: ReasonStopped, At: testTime},
		{ID: thirdID, TargetID: "bad", Kind: TargetDelivery, Reason: ReasonStopped, At: testTime},
		{ID: thirdID, TargetID: deliveryID, Kind: "bad", Reason: ReasonStopped, At: testTime},
		{ID: thirdID, TargetID: deliveryID, Kind: TargetDelivery, Reason: "bad", At: testTime},
		{ID: thirdID, TargetID: deliveryID, Kind: TargetDelivery, Reason: ReasonStopped},
	} {
		if err := value.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("tombstone=%+v err=%v", value, err)
		}
	}

	for _, location := range []TempLocation{TempLocal, TempRemote} {
		value := OwnedTemp{ID: thirdID, OwnerID: deliveryID, Location: location, Path: "/private/temp", CreatedAt: testTime}
		if err := value.Validate(); err != nil {
			t.Fatalf("temporary=%+v err=%v", value, err)
		}
	}
	for _, value := range []OwnedTemp{
		{ID: "bad", OwnerID: deliveryID, Location: TempLocal, Path: "x", CreatedAt: testTime},
		{ID: thirdID, OwnerID: "bad", Location: TempLocal, Path: "x", CreatedAt: testTime},
		{ID: thirdID, OwnerID: deliveryID, Location: "bad", Path: "x", CreatedAt: testTime},
		{ID: thirdID, OwnerID: deliveryID, Location: TempLocal, Path: "\x00", CreatedAt: testTime},
		{ID: thirdID, OwnerID: deliveryID, Location: TempLocal, Path: "x"},
	} {
		if err := value.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("temporary=%+v err=%v", value, err)
		}
	}

	for _, kind := range []HistoryKind{HistoryRegistered, HistoryActivated, HistoryStopped, HistoryFailed} {
		value := HistoryEvent{ID: thirdID, TargetID: deliveryID, Kind: kind, At: testTime, Stage: progress.StageComplete, Message: "safe failure"}
		if err := value.Validate(); err != nil {
			t.Fatalf("history=%+v err=%v", value, err)
		}
	}
	for _, value := range []HistoryEvent{
		{ID: "bad", TargetID: deliveryID, Kind: HistoryStopped, At: testTime},
		{ID: thirdID, TargetID: "bad", Kind: HistoryStopped, At: testTime},
		{ID: thirdID, TargetID: deliveryID, Kind: "bad", At: testTime},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: time.Unix(-1, 0)},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: testTime, Counters: CounterSnapshot{Sent: -1}},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: testTime, Stage: "unknown"},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: testTime, Message: "password=unsafe"},
	} {
		if err := value.Validate(); !errors.Is(err, ErrInvalid) {
			t.Fatalf("history=%+v err=%v", value, err)
		}
	}

	if invalidText("safe") || !invalidText(" ") || !invalidText("bad\x7f") || !invalidTimes(time.Time{}, testTime) || !invalidTimes(testTime, testTime.Add(-time.Second)) || invalidTimes(testTime, testTime) {
		t.Fatal("primitive validation mismatch")
	}
}

func TestCountersConcurrentAndBounded(t *testing.T) {
	if _, err := NewCounters(CounterSnapshot{Read: -1}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("new counters=%v", err)
	}
	counters, err := NewCounters(CounterSnapshot{Read: 3, Sent: 2, Confirmed: 1})
	if err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	for range 100 {
		group.Add(1)
		go func() {
			defer group.Done()
			if err := counters.AddRead(1); err != nil {
				t.Error(err)
			}
			if err := counters.AddSent(1); err != nil {
				t.Error(err)
			}
			if err := counters.AddConfirmed(1); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	if got := counters.Snapshot(); got != (CounterSnapshot{Read: 103, Sent: 102, Confirmed: 101}) {
		t.Fatalf("snapshot=%+v", got)
	}
	if err := counters.AddRead(-1); !errors.Is(err, ErrInvalid) {
		t.Fatalf("negative=%v", err)
	}
	if err := newCountersForTest().AddSent(1); !errors.Is(err, ErrInvalid) {
		t.Fatalf("counter ordering=%v", err)
	}
	overflow, err := NewCounters(CounterSnapshot{Read: math.MaxInt64})
	if err != nil {
		t.Fatal(err)
	}
	if err := overflow.AddRead(1); !errors.Is(err, ErrInvalid) {
		t.Fatalf("overflow=%v", err)
	}
}

func newCountersForTest() *Counters {
	counters, _ := NewCounters(CounterSnapshot{})
	return counters
}

func TestRegistryLifecycle(t *testing.T) {
	registry, err := NewRegistry(EmptySnapshot())
	if err != nil {
		t.Fatal(err)
	}
	server := validServer()
	item := validDelivery()
	if err := registry.RegisterServer(server); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterDelivery(item); err != nil {
		t.Fatal(err)
	}
	if err := registry.TransitionServer(thirdID, StateActive, testTime); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing server transition=%v", err)
	}
	if err := registry.TransitionServer(serverID, StateStopped, testTime.Add(time.Second)); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid server transition=%v", err)
	}
	if err := registry.TransitionServer(serverID, StateActive, time.Time{}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("server transition time=%v", err)
	}
	if err := registry.TransitionServer(serverID, StateActive, testTime.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterDelivery(item); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate active delivery=%v", err)
	}
	if err := registry.RegisterServer(server); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate server=%v", err)
	}
	duplicateDelivery := item
	duplicateDelivery.ID = server.ID
	if err := registry.RegisterDelivery(duplicateDelivery); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate delivery=%v", err)
	}
	missingDelivery := item
	missingDelivery.ID = fourthID
	missingDelivery.ServerID = thirdID
	if err := registry.RegisterDelivery(missingDelivery); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing server=%v", err)
	}
	invalidServer := server
	invalidServer.ID = "bad"
	if err := registry.RegisterServer(invalidServer); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid server=%v", err)
	}
	invalidDelivery := item
	invalidDelivery.ID = "bad"
	if err := registry.RegisterDelivery(invalidDelivery); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid delivery=%v", err)
	}

	nextPolicy := DefaultPolicy()
	nextPolicy.Version = 2
	nextPolicy.Auth = AuthBasic
	if err := registry.UpdatePolicy(deliveryID, 1, nextPolicy, testTime.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := registry.UpdatePolicy(thirdID, 1, nextPolicy, testTime); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing policy=%v", err)
	}
	if err := registry.UpdatePolicy(deliveryID, 1, nextPolicy, testTime); !errors.Is(err, ErrRevisionConflict) {
		t.Fatalf("policy conflict=%v", err)
	}
	invalidPolicy := nextPolicy
	invalidPolicy.Version = 3
	invalidPolicy.Auth = "bad"
	if err := registry.UpdatePolicy(deliveryID, 2, invalidPolicy, testTime.Add(2*time.Second)); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid policy=%v", err)
	}
	validPolicy := nextPolicy
	validPolicy.Version = 3
	if err := registry.UpdatePolicy(deliveryID, 2, validPolicy, testTime); !errors.Is(err, ErrInvalid) {
		t.Fatalf("policy time=%v", err)
	}

	if err := registry.TransitionDelivery(thirdID, StateActive, testTime); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing transition=%v", err)
	}
	if err := registry.TransitionDelivery(deliveryID, StateStopped, testTime.Add(2*time.Second)); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid transition=%v", err)
	}
	if err := registry.TransitionDelivery(deliveryID, StateActive, time.Time{}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("transition time=%v", err)
	}
	if err := registry.TransitionDelivery(deliveryID, StateActive, testTime.Add(2*time.Second)); err != nil {
		t.Fatal(err)
	}

	if err := registry.SetCounters(thirdID, CounterSnapshot{}, testTime); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing counters=%v", err)
	}
	if err := registry.SetCounters(deliveryID, CounterSnapshot{Read: 2, Sent: 2, Confirmed: 1}, testTime.Add(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := registry.SetCounters(deliveryID, CounterSnapshot{Read: 1, Sent: 2, Confirmed: 1}, testTime.Add(4*time.Second)); !errors.Is(err, ErrInvalid) {
		t.Fatalf("counter regression=%v", err)
	}

	temporary := OwnedTemp{ID: thirdID, OwnerID: deliveryID, Location: TempLocal, Path: "stage", CreatedAt: testTime}
	if err := registry.AddOwnedTemp(temporary); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddOwnedTemp(temporary); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate temp=%v", err)
	}
	missingTemp := temporary
	missingTemp.ID = fourthID
	missingTemp.OwnerID = ID("00000000-0000-4000-8000-000000000099")
	if err := registry.AddOwnedTemp(missingTemp); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing owner=%v", err)
	}
	invalidTemp := temporary
	invalidTemp.Path = ""
	if err := registry.AddOwnedTemp(invalidTemp); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid temp=%v", err)
	}
	if err := registry.RemoveOwnedTemp(fourthID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing remove=%v", err)
	}
	if err := registry.RemoveOwnedTemp(thirdID); err != nil {
		t.Fatal(err)
	}

	tombstone := Tombstone{ID: thirdID, TargetID: deliveryID, Kind: TargetDelivery, Reason: ReasonStopped, At: testTime.Add(5 * time.Second)}
	wrong := tombstone
	wrong.Kind = TargetServer
	if err := registry.TombstoneDelivery(deliveryID, wrong); !errors.Is(err, ErrInvalid) {
		t.Fatalf("wrong tombstone=%v", err)
	}
	invalidTombstone := tombstone
	invalidTombstone.ID = "bad"
	if err := registry.TombstoneDelivery(deliveryID, invalidTombstone); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid tombstone=%v", err)
	}
	duplicateTombstone := tombstone
	duplicateTombstone.ID = serverID
	if err := registry.TombstoneDelivery(deliveryID, duplicateTombstone); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate tombstone=%v", err)
	}
	if err := registry.TombstoneDelivery(thirdID, tombstone); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing tombstone=%v", err)
	}
	if err := registry.TombstoneServer(serverID, Tombstone{ID: fourthID, TargetID: serverID, Kind: TargetServer, Reason: ReasonStopped, At: testTime}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("owned delivery=%v", err)
	}
	if err := registry.TombstoneDelivery(deliveryID, tombstone); err != nil {
		t.Fatal(err)
	}
	if err := registry.AddOwnedTemp(OwnedTemp{ID: fourthID, OwnerID: deliveryID, Location: TempRemote, Path: "/stage", CreatedAt: testTime}); err != nil {
		t.Fatal(err)
	}

	serverTombstone := Tombstone{ID: ID("00000000-0000-4000-8000-000000000005"), TargetID: serverID, Kind: TargetServer, Reason: ReasonStopped, At: testTime}
	wrongServer := serverTombstone
	wrongServer.Kind = TargetDelivery
	if err := registry.TombstoneServer(serverID, wrongServer); !errors.Is(err, ErrInvalid) {
		t.Fatalf("wrong server tombstone=%v", err)
	}
	invalidServerTombstone := serverTombstone
	invalidServerTombstone.ID = "bad"
	if err := registry.TombstoneServer(serverID, invalidServerTombstone); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid server tombstone=%v", err)
	}
	duplicateServerTombstone := serverTombstone
	duplicateServerTombstone.ID = thirdID
	if err := registry.TombstoneServer(serverID, duplicateServerTombstone); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("duplicate server tombstone=%v", err)
	}
	if err := registry.TombstoneServer(ID("00000000-0000-4000-8000-000000000099"), serverTombstone); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing server tombstone=%v", err)
	}
	if err := registry.TombstoneServer(serverID, serverTombstone); err != nil {
		t.Fatal(err)
	}

	snapshot := registry.Snapshot()
	if len(snapshot.Servers) != 0 || len(snapshot.Deliveries) != 0 || len(snapshot.Tombstones) != 2 || len(snapshot.OwnedTemps) != 1 {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	snapshot.Tombstones = nil
	if len(registry.Snapshot().Tombstones) != 2 {
		t.Fatal("snapshot was not cloned")
	}
}

func TestSnapshotValidationAndCloning(t *testing.T) {
	server := validServer()
	item := validDelivery()
	temporary := OwnedTemp{ID: thirdID, OwnerID: deliveryID, Location: TempLocal, Path: "stage", CreatedAt: testTime}
	valid := Snapshot{SchemaVersion: SchemaVersion, Revision: 1, Servers: []Server{server}, Deliveries: []Delivery{item}, OwnedTemps: []OwnedTemp{temporary}}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	registry, err := NewRegistry(valid)
	if err != nil {
		t.Fatal(err)
	}
	copy := registry.Snapshot()
	copy.Deliveries[0].Policy.AllowIP = append(copy.Deliveries[0].Policy.AllowIP, "127.0.0.1/32")
	if reflect.DeepEqual(copy, registry.Snapshot()) {
		t.Fatal("policy slice was not cloned")
	}
	unsortedPolicy := DefaultPolicy()
	unsortedPolicy.AllowIP = []string{"2001:db8::/32", "127.0.0.1/32"}
	item.Policy = unsortedPolicy
	withPolicy := cloneSnapshot(Snapshot{SchemaVersion: SchemaVersion, Servers: []Server{server}, Deliveries: []Delivery{item}})
	if !reflect.DeepEqual(withPolicy.Deliveries[0].Policy.AllowIP, []string{"127.0.0.1/32", "2001:db8::/32"}) {
		t.Fatalf("allowIp order=%v", withPolicy.Deliveries[0].Policy.AllowIP)
	}

	for index, mutate := range []func(*Snapshot){
		func(value *Snapshot) { value.SchemaVersion = 2 },
		func(value *Snapshot) { value.Servers[0].ID = "bad" },
		func(value *Snapshot) { value.Deliveries[0].ID = serverID },
		func(value *Snapshot) { value.Deliveries[0].ServerID = fourthID },
		func(value *Snapshot) {
			value.Tombstones = []Tombstone{{ID: "bad", TargetID: fourthID, Kind: TargetDelivery, Reason: ReasonStopped, At: testTime}}
		},
		func(value *Snapshot) {
			value.Tombstones = []Tombstone{{ID: fourthID, TargetID: deliveryID, Kind: TargetDelivery, Reason: ReasonStopped, At: testTime}}
		},
		func(value *Snapshot) {
			value.Tombstones = []Tombstone{{ID: fourthID, TargetID: ID("00000000-0000-4000-8000-000000000005"), Kind: TargetDelivery, Reason: ReasonStopped, At: testTime}, {ID: ID("00000000-0000-4000-8000-000000000006"), TargetID: ID("00000000-0000-4000-8000-000000000005"), Kind: TargetDelivery, Reason: ReasonStale, At: testTime}}
		},
		func(value *Snapshot) { value.OwnedTemps[0].ID = serverID },
		func(value *Snapshot) { value.OwnedTemps[0].OwnerID = fourthID },
	} {
		candidate := cloneSnapshot(valid)
		mutate(&candidate)
		if err := candidate.Validate(); err == nil {
			t.Fatalf("case=%d accepted", index)
		}
		if _, err := NewRegistry(candidate); err == nil {
			t.Fatalf("new registry case=%d accepted", index)
		}
	}

	secondServer := server
	secondServer.ID = ID("00000000-0000-4000-8000-000000000010")
	secondDelivery := item
	secondDelivery.ID = ID("00000000-0000-4000-8000-000000000011")
	secondDelivery.ServerID = secondServer.ID
	secondTemp := temporary
	secondTemp.ID = ID("00000000-0000-4000-8000-000000000012")
	secondTemp.OwnerID = secondDelivery.ID
	unsorted := Snapshot{
		SchemaVersion: SchemaVersion,
		Servers:       []Server{secondServer, server},
		Deliveries:    []Delivery{secondDelivery, item},
		OwnedTemps:    []OwnedTemp{secondTemp, temporary},
	}
	cloned := cloneSnapshot(unsorted)
	if cloned.Servers[0].ID != serverID || cloned.Deliveries[0].ID != deliveryID || cloned.OwnedTemps[0].ID != thirdID {
		t.Fatalf("unsorted clone=%+v", cloned)
	}
}
