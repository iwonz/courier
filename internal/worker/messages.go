// Package worker coordinates Courier's private long-lived delivery workers.
package worker

import (
	"fmt"
	"strings"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

const MaxCompatibilityLength = 256

type HelloRequest struct {
	ServerID      delivery.ID `json:"serverId"`
	Compatibility string      `json:"compatibility"`
}

func (request HelloRequest) Validate() error {
	if !request.ServerID.Valid() || invalidCompatibility(request.Compatibility) {
		return fmt.Errorf("%w: invalid hello request", ipc.ErrProtocol)
	}
	return nil
}

type HelloResponse struct {
	ServerID      delivery.ID `json:"serverId"`
	Compatibility string      `json:"compatibility"`
}

type RegisterRequest struct {
	Delivery delivery.Delivery `json:"delivery"`
	LeaseID  delivery.ID       `json:"leaseId,omitempty"`
}

func (request RegisterRequest) Validate() error {
	if err := request.Delivery.Validate(); err != nil {
		return err
	}
	if request.LeaseID != "" && !request.LeaseID.Valid() {
		return fmt.Errorf("%w: invalid lease ID", ipc.ErrProtocol)
	}
	return nil
}

type RegisterResponse struct {
	ServerID   delivery.ID `json:"serverId"`
	DeliveryID delivery.ID `json:"deliveryId"`
	LeaseID    delivery.ID `json:"leaseId,omitempty"`
}

type LeaseRequest struct {
	DeliveryID delivery.ID `json:"deliveryId"`
	LeaseID    delivery.ID `json:"leaseId"`
}

func (request LeaseRequest) Validate() error {
	if !request.DeliveryID.Valid() || !request.LeaseID.Valid() {
		return fmt.Errorf("%w: invalid lease request", ipc.ErrProtocol)
	}
	return nil
}

type ReleaseLeaseRequest struct {
	DeliveryID delivery.ID `json:"deliveryId"`
	LeaseID    delivery.ID `json:"leaseId"`
}

func (request ReleaseLeaseRequest) Validate() error {
	return LeaseRequest(request).Validate()
}

type ListResponse struct {
	Snapshot delivery.Snapshot `json:"snapshot"`
}

type UpdatePolicyRequest struct {
	DeliveryID      delivery.ID     `json:"deliveryId"`
	ExpectedVersion uint64          `json:"expectedVersion"`
	Policy          delivery.Policy `json:"policy"`
}

func (request UpdatePolicyRequest) Validate() error {
	if !request.DeliveryID.Valid() || request.ExpectedVersion == 0 || request.Policy.Version != request.ExpectedVersion+1 {
		return fmt.Errorf("%w: invalid policy update request", ipc.ErrProtocol)
	}
	return request.Policy.Validate()
}

type TargetRequest struct {
	ID delivery.ID `json:"id"`
}

func (request TargetRequest) Validate() error {
	if !request.ID.Valid() {
		return fmt.Errorf("%w: invalid target ID", ipc.ErrProtocol)
	}
	return nil
}

type ProgressRequest struct {
	DeliveryID delivery.ID `json:"deliveryId,omitempty"`
}

func (request ProgressRequest) Validate() error {
	if request.DeliveryID != "" && !request.DeliveryID.Valid() {
		return fmt.Errorf("%w: invalid progress target", ipc.ErrProtocol)
	}
	return nil
}

type ProgressEvent struct {
	Snapshot delivery.Snapshot `json:"snapshot"`
}

func invalidCompatibility(value string) bool {
	if strings.TrimSpace(value) == "" || len(value) > MaxCompatibilityLength {
		return true
	}
	return strings.IndexFunc(value, func(character rune) bool { return character < 0x20 || character == 0x7f }) >= 0
}
