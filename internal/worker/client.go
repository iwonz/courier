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

var (
	callIPC       = ipc.Call
	openIPC       = ipc.Open
	configureIPC  = ipc.Configure
	newIPCRequest = ipc.NewRequest
)

type Client struct {
	Endpoint      string
	ServerID      delivery.ID
	Compatibility string
}

func (client Client) Validate() error {
	if client.Endpoint == "" || !client.ServerID.Valid() || invalidCompatibility(client.Compatibility) {
		return fmt.Errorf("%w: invalid worker client", ipc.ErrProtocol)
	}
	return nil
}

func (client Client) Hello(ctx context.Context) (HelloResponse, error) {
	if err := client.Validate(); err != nil {
		return HelloResponse{}, err
	}
	request, err := newIPCRequest(ipc.OperationHello, HelloRequest{ServerID: client.ServerID, Compatibility: client.Compatibility})
	if err != nil {
		return HelloResponse{}, err
	}
	response, err := callIPC(ctx, client.Endpoint, request)
	if err != nil {
		return HelloResponse{}, err
	}
	var result HelloResponse
	if err := ipc.DecodePayload(response.Payload, &result); err != nil {
		return HelloResponse{}, err
	}
	if result.ServerID != client.ServerID || result.Compatibility != client.Compatibility {
		return HelloResponse{}, fmt.Errorf("%w: hello response mismatch", ipc.ErrProtocol)
	}
	return result, nil
}

func (client Client) Register(ctx context.Context, item delivery.Delivery, foreground bool) (RegisterResponse, *Lease, error) {
	if err := client.Validate(); err != nil {
		return RegisterResponse{}, nil, err
	}
	payload := RegisterRequest{Delivery: item}
	if foreground {
		payload.LeaseID = delivery.NewID()
	}
	request, err := newIPCRequest(ipc.OperationRegister, payload)
	if err != nil {
		return RegisterResponse{}, nil, err
	}
	connection, stop, err := openIPC(ctx, client.Endpoint)
	if err != nil {
		return RegisterResponse{}, nil, err
	}
	keep := false
	defer func() {
		stop()
		if !keep {
			_ = connection.Close()
		}
	}()
	response, err := exchange(connection, request)
	if err != nil {
		return RegisterResponse{}, nil, err
	}
	var result RegisterResponse
	if err := ipc.DecodePayload(response.Payload, &result); err != nil {
		return RegisterResponse{}, nil, err
	}
	if result.ServerID != client.ServerID || result.DeliveryID != item.ID || result.LeaseID != payload.LeaseID {
		return RegisterResponse{}, nil, fmt.Errorf("%w: registration response mismatch", ipc.ErrProtocol)
	}
	if !foreground {
		return result, nil, nil
	}
	if err := connection.SetDeadline(time.Time{}); err != nil {
		return RegisterResponse{}, nil, err
	}
	keep = true
	return result, &Lease{connection: connection, deliveryID: item.ID, leaseID: payload.LeaseID}, nil
}

func (client Client) ClaimLease(ctx context.Context, deliveryID delivery.ID) (*Lease, error) {
	if err := client.Validate(); err != nil {
		return nil, err
	}
	payload := LeaseRequest{DeliveryID: deliveryID, LeaseID: delivery.NewID()}
	request, err := newIPCRequest(ipc.OperationLease, payload)
	if err != nil {
		return nil, err
	}
	connection, stop, err := openIPC(ctx, client.Endpoint)
	if err != nil {
		return nil, err
	}
	keep := false
	defer func() {
		stop()
		if !keep {
			_ = connection.Close()
		}
	}()
	if _, err := exchange(connection, request); err != nil {
		return nil, err
	}
	if err := connection.SetDeadline(time.Time{}); err != nil {
		return nil, err
	}
	keep = true
	return &Lease{connection: connection, deliveryID: deliveryID, leaseID: payload.LeaseID}, nil
}

func (client Client) List(ctx context.Context) (delivery.Snapshot, error) {
	request, err := client.request(ipc.OperationList, nil)
	if err != nil {
		return delivery.Snapshot{}, err
	}
	response, err := callIPC(ctx, client.Endpoint, request)
	if err != nil {
		return delivery.Snapshot{}, err
	}
	var result ListResponse
	if err := ipc.DecodePayload(response.Payload, &result); err != nil {
		return delivery.Snapshot{}, err
	}
	if err := result.Snapshot.Validate(); err != nil {
		return delivery.Snapshot{}, err
	}
	return result.Snapshot, nil
}

func (client Client) UpdatePolicy(ctx context.Context, update UpdatePolicyRequest) error {
	return client.call(ctx, ipc.OperationUpdatePolicy, update)
}

func (client Client) StopDelivery(ctx context.Context, id delivery.ID) error {
	return client.call(ctx, ipc.OperationStopDelivery, TargetRequest{ID: id})
}

func (client Client) StopServer(ctx context.Context) error {
	return client.call(ctx, ipc.OperationStopServer, nil)
}

func (client Client) Shutdown(ctx context.Context) error {
	return client.call(ctx, ipc.OperationShutdown, nil)
}

func (client Client) Subscribe(ctx context.Context, id delivery.ID) (*Subscription, error) {
	request, err := client.request(ipc.OperationSubscribeProgress, ProgressRequest{DeliveryID: id})
	if err != nil {
		return nil, err
	}
	connection, stop, err := openIPC(ctx, client.Endpoint)
	if err != nil {
		return nil, err
	}
	if err := ipc.WriteFrame(connection, request); err != nil {
		stop()
		_ = connection.Close()
		return nil, err
	}
	stop()
	if err := connection.SetDeadline(time.Time{}); err != nil {
		_ = connection.Close()
		return nil, err
	}
	return &Subscription{connection: connection, requestID: request.ID}, nil
}

func (client Client) request(operation ipc.Operation, payload any) (ipc.Request, error) {
	if err := client.Validate(); err != nil {
		return ipc.Request{}, err
	}
	return newIPCRequest(operation, payload)
}

func (client Client) call(ctx context.Context, operation ipc.Operation, payload any) error {
	request, err := client.request(operation, payload)
	if err != nil {
		return err
	}
	_, err = callIPC(ctx, client.Endpoint, request)
	return err
}

type Lease struct {
	connection net.Conn
	deliveryID delivery.ID
	leaseID    delivery.ID
	once       sync.Once
	err        error
}

func (lease *Lease) ID() delivery.ID { return lease.leaseID }

func (lease *Lease) Release(ctx context.Context) error {
	lease.once.Do(func() {
		stop, err := configureIPC(ctx, lease.connection)
		if err == nil {
			request, requestErr := newIPCRequest(ipc.OperationReleaseLease, ReleaseLeaseRequest{DeliveryID: lease.deliveryID, LeaseID: lease.leaseID})
			if requestErr == nil {
				_, err = exchange(lease.connection, request)
			} else {
				err = requestErr
			}
			stop()
		}
		lease.err = errors.Join(err, lease.connection.Close())
	})
	return lease.err
}

func (lease *Lease) Close() error {
	lease.once.Do(func() { lease.err = lease.connection.Close() })
	return lease.err
}

type Subscription struct {
	connection net.Conn
	requestID  delivery.ID
	once       sync.Once
	err        error
}

func (subscription *Subscription) Next(ctx context.Context) (ProgressEvent, error) {
	stop, err := configureIPC(ctx, subscription.connection)
	if err != nil {
		return ProgressEvent{}, err
	}
	defer stop()
	var response ipc.Response
	if err := ipc.ReadFrame(subscription.connection, &response); err != nil {
		return ProgressEvent{}, err
	}
	if err := response.Validate(); err != nil {
		return ProgressEvent{}, err
	}
	if response.ID != subscription.requestID {
		return ProgressEvent{}, fmt.Errorf("%w: subscription response mismatch", ipc.ErrProtocol)
	}
	if response.Error != nil {
		return ProgressEvent{}, &ipc.RemoteError{Code: response.Error.Code, Message: response.Error.Message}
	}
	var event ProgressEvent
	if err := ipc.DecodePayload(response.Payload, &event); err != nil {
		return ProgressEvent{}, err
	}
	if err := event.Snapshot.Validate(); err != nil {
		return ProgressEvent{}, err
	}
	return event, nil
}

func (subscription *Subscription) Close() error {
	subscription.once.Do(func() { subscription.err = subscription.connection.Close() })
	return subscription.err
}

func exchange(connection net.Conn, request ipc.Request) (ipc.Response, error) {
	if err := ipc.WriteFrame(connection, request); err != nil {
		return ipc.Response{}, err
	}
	var response ipc.Response
	if err := ipc.ReadFrame(connection, &response); err != nil {
		return ipc.Response{}, err
	}
	if err := response.Validate(); err != nil {
		return ipc.Response{}, err
	}
	if response.ID != request.ID {
		return ipc.Response{}, fmt.Errorf("%w: response ID mismatch", ipc.ErrProtocol)
	}
	if response.Error != nil {
		return response, &ipc.RemoteError{Code: response.Error.Code, Message: response.Error.Message}
	}
	return response, nil
}
