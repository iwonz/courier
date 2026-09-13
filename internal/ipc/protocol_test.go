package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/delivery"
)

const requestID delivery.ID = "00000000-0000-4000-8000-000000000101"

func validRequest() Request {
	return Request{Version: ProtocolVersion, ID: requestID, Operation: OperationHello, Payload: json.RawMessage(`{"value":"ok"}`)}
}

func TestProtocolEnumsAndRequestValidation(t *testing.T) {
	operations := []Operation{
		OperationHello, OperationRegister, OperationLease, OperationReleaseLease, OperationList,
		OperationUpdatePolicy, OperationStopDelivery, OperationStopServer, OperationSubscribeProgress, OperationShutdown,
	}
	for _, operation := range operations {
		if !operation.Valid() {
			t.Fatalf("operation=%q", operation)
		}
	}
	if Operation("unknown").Valid() {
		t.Fatal("unknown operation accepted")
	}
	for _, code := range []ErrorCode{CodeInvalid, CodeVersion, CodeUnsupported, CodeNotFound, CodeConflict, CodeUnavailable, CodeInternal} {
		if !code.Valid() {
			t.Fatalf("code=%q", code)
		}
	}
	if ErrorCode("unknown").Valid() {
		t.Fatal("unknown error code accepted")
	}

	request, err := NewRequest(OperationHello, map[string]string{"value": "ok"})
	if err != nil || !request.ID.Valid() || string(request.Payload) != `{"value":"ok"}` {
		t.Fatalf("request=%+v err=%v", request, err)
	}
	empty, err := NewRequest(OperationList, nil)
	if err != nil || len(empty.Payload) != 0 {
		t.Fatalf("empty=%+v err=%v", empty, err)
	}
	if _, err := NewRequest(OperationHello, make(chan int)); err == nil {
		t.Fatal("unsupported request payload accepted")
	}
	if _, err := NewRequest("unknown", nil); !errors.Is(err, ErrProtocol) {
		t.Fatalf("unknown operation=%v", err)
	}

	for _, mutate := range []func(*Request){
		func(value *Request) { value.Version = 2 },
		func(value *Request) { value.ID = "bad" },
		func(value *Request) { value.Operation = "bad" },
		func(value *Request) { value.Payload = json.RawMessage(`{`) },
	} {
		candidate := validRequest()
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, ErrProtocol) {
			t.Fatalf("request=%+v err=%v", candidate, err)
		}
	}
}

func TestResponsesAndPayloads(t *testing.T) {
	request := validRequest()
	success, err := Success(request, map[string]bool{"ready": true})
	if err != nil || !success.OK || string(success.Payload) != `{"ready":true}` {
		t.Fatalf("success=%+v err=%v", success, err)
	}
	empty, err := Success(request, nil)
	if err != nil || len(empty.Payload) != 0 {
		t.Fatalf("empty=%+v err=%v", empty, err)
	}
	if _, err := Success(request, make(chan int)); err == nil {
		t.Fatal("unsupported response payload accepted")
	}
	rejection, err := Rejection(request, CodeConflict, "already registered")
	if err != nil || rejection.OK || rejection.Error.Code != CodeConflict {
		t.Fatalf("rejection=%+v err=%v", rejection, err)
	}
	if _, err := Rejection(request, "bad", "failure"); !errors.Is(err, ErrProtocol) {
		t.Fatalf("bad rejection=%v", err)
	}
	if _, err := Rejection(request, CodeInvalid, "bad\nmessage"); !errors.Is(err, ErrProtocol) {
		t.Fatalf("bad message=%v", err)
	}

	invalid := []Response{
		{Version: 2, ID: requestID, OK: true},
		{Version: ProtocolVersion, ID: "bad", OK: true},
		{Version: ProtocolVersion, ID: requestID, OK: true, Error: &Failure{Code: CodeInvalid, Message: "bad"}},
		{Version: ProtocolVersion, ID: requestID},
		{Version: ProtocolVersion, ID: requestID, Error: &Failure{Code: "bad", Message: "bad"}},
		{Version: ProtocolVersion, ID: requestID, Error: &Failure{Code: CodeInvalid}},
		{Version: ProtocolVersion, ID: requestID, Error: &Failure{Code: CodeInvalid, Message: strings.Repeat("x", 1025)}},
		{Version: ProtocolVersion, ID: requestID, OK: true, Payload: json.RawMessage(`{`)},
	}
	for _, response := range invalid {
		if err := response.Validate(); !errors.Is(err, ErrProtocol) {
			t.Fatalf("response=%+v err=%v", response, err)
		}
	}

	var payload struct {
		Ready bool `json:"ready"`
	}
	if err := DecodePayload(success.Payload, &payload); err != nil || !payload.Ready {
		t.Fatalf("payload=%+v err=%v", payload, err)
	}
	if err := DecodePayload(success.Payload, nil); !errors.Is(err, ErrProtocol) {
		t.Fatalf("nil payload target=%v", err)
	}
	if err := DecodePayload(nil, &payload); !errors.Is(err, ErrProtocol) {
		t.Fatalf("missing payload=%v", err)
	}
	if err := DecodePayload(json.RawMessage(`{"ready":true,"extra":1}`), &payload); !errors.Is(err, ErrProtocol) {
		t.Fatalf("unknown payload field=%v", err)
	}
	if err := DecodePayload(json.RawMessage(`{"ready":true}{}`), &payload); !errors.Is(err, ErrProtocol) {
		t.Fatalf("trailing payload=%v", err)
	}
}

func TestFrameCodec(t *testing.T) {
	request := validRequest()
	var encoded bytes.Buffer
	if err := WriteFrame(&encoded, request); err != nil {
		t.Fatal(err)
	}
	var decoded Request
	if err := ReadFrame(&encoded, &decoded); err != nil || decoded.ID != request.ID {
		t.Fatalf("decoded=%+v err=%v", decoded, err)
	}
	if err := ReadFrame(bytes.NewReader(nil), &decoded); !errors.Is(err, io.EOF) {
		t.Fatalf("short header=%v", err)
	}
	if err := ReadFrame(bytes.NewReader([]byte{0, 0, 0, 0}), &decoded); !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("empty frame=%v", err)
	}
	largeHeader := [4]byte{}
	binary.BigEndian.PutUint32(largeHeader[:], MaxFrameSize+1)
	if err := ReadFrame(bytes.NewReader(largeHeader[:]), &decoded); !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("large frame=%v", err)
	}
	short := framed([]byte(`{"version":1}`))
	short = short[:len(short)-1]
	if err := ReadFrame(bytes.NewReader(short), &decoded); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("short body=%v", err)
	}
	if err := ReadFrame(bytes.NewReader(framed([]byte(`{"unknown":true}`))), &decoded); !errors.Is(err, ErrProtocol) {
		t.Fatalf("invalid JSON envelope=%v", err)
	}
	if err := ReadFrame(bytes.NewReader(framed([]byte(`{}{} `))), &decoded); !errors.Is(err, ErrProtocol) {
		t.Fatalf("trailing JSON=%v", err)
	}
	if err := ReadFrame(bytes.NewReader(framed([]byte(`{}`))), nil); !errors.Is(err, ErrProtocol) {
		t.Fatalf("nil destination=%v", err)
	}

	if err := WriteFrame(io.Discard, make(chan int)); err == nil {
		t.Fatal("unsupported frame value accepted")
	}
	if err := WriteFrame(io.Discard, strings.Repeat("x", MaxFrameSize)); !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("large write=%v", err)
	}
	writer := &stepWriter{failAt: 1, err: errors.New("header")}
	if err := WriteFrame(writer, request); !errors.Is(err, writer.err) {
		t.Fatalf("header write=%v", err)
	}
	writer = &stepWriter{failAt: 2, err: errors.New("body")}
	if err := WriteFrame(writer, request); !errors.Is(err, writer.err) {
		t.Fatalf("body write=%v", err)
	}
	writer = &stepWriter{zeroAt: 1}
	if err := WriteFrame(writer, request); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("zero write=%v", err)
	}
	writer = &stepWriter{oversizeAt: 1}
	if err := WriteFrame(writer, request); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("oversize write=%v", err)
	}
	writer = &stepWriter{chunk: 2}
	if err := WriteFrame(writer, request); err != nil || writer.Buffer.Len() == 0 {
		t.Fatalf("partial writes=%v", err)
	}
}

func framed(data []byte) []byte {
	result := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(result, uint32(len(data)))
	copy(result[4:], data)
	return result
}

type stepWriter struct {
	bytes.Buffer
	calls      int
	failAt     int
	zeroAt     int
	oversizeAt int
	chunk      int
	err        error
}

func (writer *stepWriter) Write(data []byte) (int, error) {
	writer.calls++
	if writer.calls == writer.failAt {
		return 0, writer.err
	}
	if writer.calls == writer.zeroAt {
		return 0, nil
	}
	if writer.calls == writer.oversizeAt {
		return len(data) + 1, nil
	}
	if writer.chunk > 0 && len(data) > writer.chunk {
		data = data[:writer.chunk]
	}
	return writer.Buffer.Write(data)
}
