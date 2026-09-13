package worker

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

const (
	workerServerID   delivery.ID = "00000000-0000-4000-8000-000000000201"
	workerDeliveryID delivery.ID = "00000000-0000-4000-8000-000000000202"
	workerLeaseID    delivery.ID = "00000000-0000-4000-8000-000000000203"
)

var workerTime = time.Unix(1_800_000_000, 0).UTC()

func validWorkerDelivery() delivery.Delivery {
	return delivery.Delivery{
		ID: workerDeliveryID, ServerID: workerServerID, Route: delivery.RoutePathToWeb, State: delivery.StateStarting,
		Policy: delivery.DefaultPolicy(), CreatedAt: workerTime, UpdatedAt: workerTime,
	}
}

func TestWorkerMessageValidation(t *testing.T) {
	hello := HelloRequest{ServerID: workerServerID, Compatibility: "http-v1"}
	if err := hello.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, value := range []HelloRequest{{ServerID: "bad", Compatibility: "http-v1"}, {ServerID: workerServerID}, {ServerID: workerServerID, Compatibility: "bad\nvalue"}, {ServerID: workerServerID, Compatibility: strings.Repeat("x", MaxCompatibilityLength+1)}} {
		if err := value.Validate(); !errors.Is(err, ipc.ErrProtocol) {
			t.Fatalf("hello=%+v err=%v", value, err)
		}
	}

	registration := RegisterRequest{Delivery: validWorkerDelivery()}
	if err := registration.Validate(); err != nil {
		t.Fatal(err)
	}
	registration.LeaseID = workerLeaseID
	if err := registration.Validate(); err != nil {
		t.Fatal(err)
	}
	registration.LeaseID = "bad"
	if err := registration.Validate(); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("lease=%v", err)
	}
	registration = RegisterRequest{Delivery: validWorkerDelivery()}
	registration.Delivery.ID = "bad"
	if err := registration.Validate(); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("delivery=%v", err)
	}

	lease := LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: workerLeaseID}
	if err := lease.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (LeaseRequest{DeliveryID: "bad", LeaseID: workerLeaseID}).Validate(); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("delivery lease=%v", err)
	}
	if err := (LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: "bad"}).Validate(); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("lease ID=%v", err)
	}
	if err := (ReleaseLeaseRequest(lease)).Validate(); err != nil {
		t.Fatal(err)
	}

	policy := delivery.DefaultPolicy()
	policy.Version = 2
	update := UpdatePolicyRequest{DeliveryID: workerDeliveryID, ExpectedVersion: 1, Policy: policy}
	if err := update.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*UpdatePolicyRequest){
		func(value *UpdatePolicyRequest) { value.DeliveryID = "bad" },
		func(value *UpdatePolicyRequest) { value.ExpectedVersion = 0 },
		func(value *UpdatePolicyRequest) { value.Policy.Version = 9 },
		func(value *UpdatePolicyRequest) { value.Policy.Auth = "bad" },
	} {
		candidate := update
		mutate(&candidate)
		if err := candidate.Validate(); err == nil {
			t.Fatalf("update=%+v accepted", candidate)
		}
	}

	if err := (TargetRequest{ID: workerDeliveryID}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (TargetRequest{ID: "bad"}).Validate(); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("target=%v", err)
	}
	if err := (ProgressRequest{}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (ProgressRequest{DeliveryID: workerDeliveryID}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (ProgressRequest{DeliveryID: "bad"}).Validate(); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("progress=%v", err)
	}
	if invalidCompatibility("http-v1") || !invalidCompatibility(" ") || !invalidCompatibility("bad\x7f") {
		t.Fatal("compatibility validation mismatch")
	}
}
