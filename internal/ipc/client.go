package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

const defaultCallTimeout = 5 * time.Second

var dialControl = Dial

type RemoteError struct {
	Code    ErrorCode
	Message string
}

func (failure *RemoteError) Error() string {
	return fmt.Sprintf("Courier IPC %s: %s", failure.Code, failure.Message)
}

func Call(ctx context.Context, endpoint string, request Request) (Response, error) {
	if err := request.Validate(); err != nil {
		return Response{}, err
	}
	connection, stop, err := Open(ctx, endpoint)
	if err != nil {
		return Response{}, err
	}
	defer connection.Close()
	defer stop()
	if err := WriteFrame(connection, request); err != nil {
		return Response{}, err
	}
	var response Response
	if err := ReadFrame(connection, &response); err != nil {
		return Response{}, err
	}
	if err := response.Validate(); err != nil {
		return Response{}, err
	}
	if response.ID != request.ID {
		return Response{}, fmt.Errorf("%w: response ID mismatch", ErrProtocol)
	}
	if response.Error != nil {
		return response, &RemoteError{Code: response.Error.Code, Message: response.Error.Message}
	}
	return response, nil
}

func Open(ctx context.Context, endpoint string) (net.Conn, func() bool, error) {
	connection, err := dialControl(ctx, endpoint)
	if err != nil {
		return nil, nil, err
	}
	stop, err := configureConnection(ctx, connection)
	if err != nil {
		_ = connection.Close()
		return nil, nil, err
	}
	return connection, stop, nil
}

func Configure(ctx context.Context, connection net.Conn) (func() bool, error) {
	if connection == nil {
		return nil, fmt.Errorf("%w: connection is required", ErrProtocol)
	}
	return configureConnection(ctx, connection)
}

func configureConnection(ctx context.Context, connection net.Conn) (func() bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(defaultCallTimeout)
	}
	if err := connection.SetDeadline(deadline); err != nil {
		return nil, err
	}
	stop := context.AfterFunc(ctx, func() { _ = connection.SetDeadline(time.Now()) })
	return stop, nil
}

func IsRemoteError(err error, code ErrorCode) bool {
	var failure *RemoteError
	return errors.As(err, &failure) && failure.Code == code
}
