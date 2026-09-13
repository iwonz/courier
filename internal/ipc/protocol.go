// Package ipc defines Courier's private, versioned local control protocol.
package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/iwonz/courier/internal/delivery"
)

const (
	ProtocolVersion = 1
	MaxFrameSize    = 1 << 20
)

var (
	ErrProtocol      = errors.New("invalid Courier IPC protocol")
	ErrFrameTooLarge = errors.New("Courier IPC frame exceeds limit")
)

type Operation string

const (
	OperationHello             Operation = "hello"
	OperationRegister          Operation = "register"
	OperationLease             Operation = "lease"
	OperationReleaseLease      Operation = "release-lease"
	OperationList              Operation = "list"
	OperationUpdatePolicy      Operation = "update-policy"
	OperationStopDelivery      Operation = "stop-delivery"
	OperationStopServer        Operation = "stop-server"
	OperationSubscribeProgress Operation = "subscribe-progress"
	OperationShutdown          Operation = "shutdown"
)

func (operation Operation) Valid() bool {
	switch operation {
	case OperationHello, OperationRegister, OperationLease, OperationReleaseLease, OperationList, OperationUpdatePolicy, OperationStopDelivery, OperationStopServer, OperationSubscribeProgress, OperationShutdown:
		return true
	default:
		return false
	}
}

type ErrorCode string

const (
	CodeInvalid     ErrorCode = "invalid"
	CodeVersion     ErrorCode = "version"
	CodeUnsupported ErrorCode = "unsupported"
	CodeNotFound    ErrorCode = "not-found"
	CodeConflict    ErrorCode = "conflict"
	CodeUnavailable ErrorCode = "unavailable"
	CodeInternal    ErrorCode = "internal"
)

func (code ErrorCode) Valid() bool {
	switch code {
	case CodeInvalid, CodeVersion, CodeUnsupported, CodeNotFound, CodeConflict, CodeUnavailable, CodeInternal:
		return true
	default:
		return false
	}
}

type Request struct {
	Version   int             `json:"version"`
	ID        delivery.ID     `json:"id"`
	Operation Operation       `json:"operation"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func NewRequest(operation Operation, payload any) (Request, error) {
	request := Request{Version: ProtocolVersion, ID: delivery.NewID(), Operation: operation}
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return Request{}, err
		}
		request.Payload = data
	}
	if err := request.Validate(); err != nil {
		return Request{}, err
	}
	return request, nil
}

func (request Request) Validate() error {
	if request.Version != ProtocolVersion {
		return fmt.Errorf("%w: unsupported version %d", ErrProtocol, request.Version)
	}
	if !request.ID.Valid() {
		return fmt.Errorf("%w: request ID must be a canonical UUID", ErrProtocol)
	}
	if !request.Operation.Valid() {
		return fmt.Errorf("%w: unknown operation %q", ErrProtocol, request.Operation)
	}
	if len(request.Payload) != 0 && !json.Valid(request.Payload) {
		return fmt.Errorf("%w: malformed payload", ErrProtocol)
	}
	return nil
}

type Failure struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type Response struct {
	Version int             `json:"version"`
	ID      delivery.ID     `json:"id"`
	OK      bool            `json:"ok"`
	Error   *Failure        `json:"error,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func Success(request Request, payload any) (Response, error) {
	response := Response{Version: ProtocolVersion, ID: request.ID, OK: true}
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return Response{}, err
		}
		response.Payload = data
	}
	return response, response.Validate()
}

func Rejection(request Request, code ErrorCode, message string) (Response, error) {
	response := Response{Version: ProtocolVersion, ID: request.ID, Error: &Failure{Code: code, Message: message}}
	return response, response.Validate()
}

func (response Response) Validate() error {
	if response.Version != ProtocolVersion || !response.ID.Valid() {
		return fmt.Errorf("%w: invalid response envelope", ErrProtocol)
	}
	if response.OK == (response.Error != nil) {
		return fmt.Errorf("%w: response result is inconsistent", ErrProtocol)
	}
	if response.Error != nil && (!response.Error.Code.Valid() || invalidMessage(response.Error.Message)) {
		return fmt.Errorf("%w: invalid response error", ErrProtocol)
	}
	if len(response.Payload) != 0 && !json.Valid(response.Payload) {
		return fmt.Errorf("%w: malformed response payload", ErrProtocol)
	}
	return nil
}

func DecodePayload(payload json.RawMessage, destination any) error {
	if destination == nil {
		return fmt.Errorf("%w: payload destination is required", ErrProtocol)
	}
	if len(payload) == 0 {
		return fmt.Errorf("%w: payload is required", ErrProtocol)
	}
	return decodeStrict(payload, destination)
}

func WriteFrame(writer io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) == 0 || len(data) > MaxFrameSize {
		return ErrFrameTooLarge
	}
	header := [4]byte{}
	binary.BigEndian.PutUint32(header[:], uint32(len(data)))
	if err := writeAll(writer, header[:]); err != nil {
		return err
	}
	return writeAll(writer, data)
}

func ReadFrame(reader io.Reader, destination any) error {
	if destination == nil {
		return fmt.Errorf("%w: frame destination is required", ErrProtocol)
	}
	header := [4]byte{}
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return err
	}
	size := binary.BigEndian.Uint32(header[:])
	if size == 0 || size > MaxFrameSize {
		return ErrFrameTooLarge
	}
	data := make([]byte, int(size))
	if _, err := io.ReadFull(reader, data); err != nil {
		return err
	}
	return decodeStrict(data, destination)
}

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("%w: %v", ErrProtocol, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: trailing JSON data", ErrProtocol)
	}
	return nil
}

func writeAll(writer io.Writer, data []byte) error {
	for len(data) != 0 {
		written, err := writer.Write(data)
		if err != nil {
			return err
		}
		if written <= 0 || written > len(data) {
			return io.ErrShortWrite
		}
		data = data[written:]
	}
	return nil
}

func invalidMessage(value string) bool {
	if value == "" || len(value) > 1024 {
		return true
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return true
		}
	}
	return false
}
