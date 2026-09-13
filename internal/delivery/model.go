// Package delivery defines the secret-free domain and durable registry used by
// Courier's long-lived HTTP deliveries and control plane.
package delivery

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const SchemaVersion = 1

var (
	ErrInvalid          = errors.New("invalid delivery registry data")
	ErrDuplicate        = errors.New("duplicate delivery registry identity")
	ErrNotFound         = errors.New("delivery registry identity not found")
	ErrRevisionConflict = errors.New("delivery registry revision conflict")
)

type ID string

func NewID() ID { return ID(uuid.NewString()) }

func ParseID(value string) (ID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil || parsed == uuid.Nil || parsed.String() != value {
		return "", fmt.Errorf("%w: expected canonical UUID", ErrInvalid)
	}
	return ID(parsed.String()), nil
}

func (id ID) Valid() bool {
	_, err := ParseID(string(id))
	return err == nil
}

type State string

const (
	StateStarting State = "starting"
	StateActive   State = "active"
	StateStopping State = "stopping"
	StateStopped  State = "stopped"
	StateFailed   State = "failed"
)

func (state State) Valid() bool {
	switch state {
	case StateStarting, StateActive, StateStopping, StateStopped, StateFailed:
		return true
	default:
		return false
	}
}

func CanTransition(from, to State) bool {
	switch from {
	case StateStarting:
		return to == StateActive || to == StateStopping || to == StateFailed
	case StateActive:
		return to == StateStopping || to == StateFailed
	case StateStopping:
		return to == StateStopped || to == StateFailed
	case StateFailed:
		return to == StateStopped
	default:
		return false
	}
}

type AuthMode string

const (
	AuthNone     AuthMode = "none"
	AuthBasic    AuthMode = "basic"
	AuthPassword AuthMode = "password"
)

type AuthFailAction string

const (
	AuthFailBan  AuthFailAction = "ban"
	AuthFailStop AuthFailAction = "stop"
)

type Limit struct {
	Unlimited bool  `json:"unlimited"`
	Value     int64 `json:"value"`
}

func (limit Limit) validate(name string) error {
	if limit.Value < 0 || limit.Unlimited && limit.Value != 0 {
		return fmt.Errorf("%w: %s limit is inconsistent", ErrInvalid, name)
	}
	return nil
}

type Policy struct {
	Version          uint64         `json:"version"`
	Auth             AuthMode       `json:"auth"`
	AuthAttempts     uint64         `json:"authAttempts"`
	AuthFailAction   AuthFailAction `json:"authFailAction"`
	DeliveryLimit    Limit          `json:"deliveryLimit"`
	AllowIP          []string       `json:"allowIp"`
	MaxFileSize      Limit          `json:"maxFileSize"`
	MaxExtractedSize Limit          `json:"maxExtractedSize"`
	UploadRate       Limit          `json:"uploadRate"`
	DownloadRate     Limit          `json:"downloadRate"`
	NoUI             bool           `json:"noUi"`
}

func DefaultPolicy() Policy {
	return Policy{
		Version:          1,
		Auth:             AuthNone,
		AuthAttempts:     5,
		AuthFailAction:   AuthFailBan,
		DeliveryLimit:    Limit{Unlimited: true},
		MaxFileSize:      Limit{Value: 10 << 30},
		MaxExtractedSize: Limit{Value: 100 << 30},
		UploadRate:       Limit{Unlimited: true},
		DownloadRate:     Limit{Unlimited: true},
	}
}

func (policy Policy) Validate() error {
	if policy.Version == 0 || policy.AuthAttempts == 0 {
		return fmt.Errorf("%w: policy version and auth attempts must be positive", ErrInvalid)
	}
	if policy.Auth != AuthNone && policy.Auth != AuthBasic && policy.Auth != AuthPassword {
		return fmt.Errorf("%w: unknown authentication mode", ErrInvalid)
	}
	if policy.AuthFailAction != AuthFailBan && policy.AuthFailAction != AuthFailStop {
		return fmt.Errorf("%w: unknown authentication failure action", ErrInvalid)
	}
	for name, limit := range map[string]Limit{
		"delivery": policy.DeliveryLimit, "file size": policy.MaxFileSize, "extracted size": policy.MaxExtractedSize,
		"upload rate": policy.UploadRate, "download rate": policy.DownloadRate,
	} {
		if err := limit.validate(name); err != nil {
			return err
		}
	}
	seen := map[netip.Prefix]bool{}
	for _, value := range policy.AllowIP {
		prefix, err := netip.ParsePrefix(value)
		if err != nil || prefix.Masked().String() != value || seen[prefix] {
			return fmt.Errorf("%w: allowIp must contain unique canonical CIDRs", ErrInvalid)
		}
		seen[prefix] = true
	}
	return nil
}

type CounterSnapshot struct {
	Read      int64 `json:"read"`
	Sent      int64 `json:"sent"`
	Confirmed int64 `json:"confirmed"`
}

func (counter CounterSnapshot) Validate() error {
	if counter.Read < 0 || counter.Sent < 0 || counter.Confirmed < 0 {
		return fmt.Errorf("%w: counters must be non-negative", ErrInvalid)
	}
	return nil
}

type Server struct {
	ID              ID        `json:"id"`
	Bind            string    `json:"bind"`
	ControlEndpoint string    `json:"controlEndpoint"`
	Compatibility   string    `json:"compatibility,omitempty"`
	ProcessID       int       `json:"processId"`
	State           State     `json:"state"`
	StartedAt       time.Time `json:"startedAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (server Server) Validate() error {
	if !server.ID.Valid() || invalidText(server.Bind) || invalidText(server.ControlEndpoint) || server.ProcessID <= 0 || !server.State.Valid() || invalidTimes(server.StartedAt, server.UpdatedAt) {
		return fmt.Errorf("%w: invalid server", ErrInvalid)
	}
	if server.Compatibility != "" && (invalidText(server.Compatibility) || len(server.Compatibility) > 256) {
		return fmt.Errorf("%w: invalid server compatibility", ErrInvalid)
	}
	_, port, err := net.SplitHostPort(server.Bind)
	number, numberErr := strconv.ParseUint(port, 10, 16)
	if err != nil || numberErr != nil || number == 0 {
		return fmt.Errorf("%w: server bind must use host:port", ErrInvalid)
	}
	return nil
}

type Route string

const (
	RouteWebToPath     Route = "web-to-path"
	RoutePathToWeb     Route = "path-to-web"
	RouteWebhookToPath Route = "webhook-to-path"
	RoutePathToHTTP    Route = "path-to-http"
)

func (route Route) Valid() bool {
	return route == RouteWebToPath || route == RoutePathToWeb || route == RouteWebhookToPath || route == RoutePathToHTTP
}

type Delivery struct {
	ID          ID              `json:"id"`
	ServerID    ID              `json:"serverId"`
	Route       Route           `json:"route"`
	Source      string          `json:"source,omitempty"`
	Destination string          `json:"destination,omitempty"`
	State       State           `json:"state"`
	Policy      Policy          `json:"policy"`
	Counters    CounterSnapshot `json:"counters"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func (delivery Delivery) Validate() error {
	if !delivery.ID.Valid() || !delivery.ServerID.Valid() || !delivery.Route.Valid() || !delivery.State.Valid() || invalidTimes(delivery.CreatedAt, delivery.UpdatedAt) {
		return fmt.Errorf("%w: invalid delivery", ErrInvalid)
	}
	if delivery.Source != "" && invalidText(delivery.Source) || delivery.Destination != "" && invalidText(delivery.Destination) {
		return fmt.Errorf("%w: invalid delivery endpoint metadata", ErrInvalid)
	}
	return errors.Join(delivery.Policy.Validate(), delivery.Counters.Validate())
}

type TargetKind string

const (
	TargetServer   TargetKind = "server"
	TargetDelivery TargetKind = "delivery"
)

type TombstoneReason string

const (
	ReasonStopped   TombstoneReason = "stopped"
	ReasonCompleted TombstoneReason = "completed"
	ReasonFailed    TombstoneReason = "failed"
	ReasonStale     TombstoneReason = "stale"
)

type Tombstone struct {
	ID       ID              `json:"id"`
	TargetID ID              `json:"targetId"`
	Kind     TargetKind      `json:"kind"`
	Reason   TombstoneReason `json:"reason"`
	At       time.Time       `json:"at"`
}

func (tombstone Tombstone) Validate() error {
	validKind := tombstone.Kind == TargetServer || tombstone.Kind == TargetDelivery
	validReason := tombstone.Reason == ReasonStopped || tombstone.Reason == ReasonCompleted || tombstone.Reason == ReasonFailed || tombstone.Reason == ReasonStale
	if !tombstone.ID.Valid() || !tombstone.TargetID.Valid() || !validKind || !validReason || tombstone.At.IsZero() {
		return fmt.Errorf("%w: invalid tombstone", ErrInvalid)
	}
	return nil
}

type TempLocation string

const (
	TempLocal  TempLocation = "local"
	TempRemote TempLocation = "remote"
)

type OwnedTemp struct {
	ID        ID           `json:"id"`
	OwnerID   ID           `json:"ownerId"`
	Location  TempLocation `json:"location"`
	Path      string       `json:"path"`
	CreatedAt time.Time    `json:"createdAt"`
}

func (temporary OwnedTemp) Validate() error {
	if !temporary.ID.Valid() || !temporary.OwnerID.Valid() || temporary.Location != TempLocal && temporary.Location != TempRemote || invalidText(temporary.Path) || temporary.CreatedAt.IsZero() {
		return fmt.Errorf("%w: invalid owned temporary resource", ErrInvalid)
	}
	return nil
}

type HistoryKind string

const (
	HistoryRegistered HistoryKind = "registered"
	HistoryActivated  HistoryKind = "activated"
	HistoryStopped    HistoryKind = "stopped"
	HistoryFailed     HistoryKind = "failed"
)

type HistoryEvent struct {
	ID       ID              `json:"id"`
	TargetID ID              `json:"targetId"`
	Kind     HistoryKind     `json:"kind"`
	At       time.Time       `json:"at"`
	Counters CounterSnapshot `json:"counters"`
}

func (event HistoryEvent) Validate() error {
	validKind := event.Kind == HistoryRegistered || event.Kind == HistoryActivated || event.Kind == HistoryStopped || event.Kind == HistoryFailed
	if !event.ID.Valid() || !event.TargetID.Valid() || !validKind || event.At.IsZero() || event.At.Before(time.Unix(0, 0)) {
		return fmt.Errorf("%w: invalid history event", ErrInvalid)
	}
	return event.Counters.Validate()
}

func invalidText(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}
	return strings.IndexFunc(value, func(r rune) bool { return r < 0x20 || r == 0x7f }) >= 0
}

func invalidTimes(created, updated time.Time) bool {
	return created.IsZero() || updated.Before(created)
}
