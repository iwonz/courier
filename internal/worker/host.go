package worker

import (
	"context"
	"encoding/json"
	"net"

	"github.com/iwonz/courier/internal/delivery"
)

// DeliveryHost owns ephemeral data-plane definitions and resources. Durable
// registry records deliberately never contain this material.
type DeliveryHost interface {
	Register(context.Context, delivery.Delivery, json.RawMessage) error
	OwnedTemps(delivery.ID) []delivery.OwnedTemp
	PreparePolicy(context.Context, delivery.ID, uint64, delivery.Policy) (commit func(), cancel func(), err error)
	Stop(context.Context, delivery.ID) error
	Serve(net.Listener) error
	Close(context.Context) error
}
