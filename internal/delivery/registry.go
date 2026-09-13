package delivery

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type Snapshot struct {
	SchemaVersion int         `json:"schemaVersion"`
	Revision      uint64      `json:"revision"`
	Servers       []Server    `json:"servers"`
	Deliveries    []Delivery  `json:"deliveries"`
	Tombstones    []Tombstone `json:"tombstones"`
	OwnedTemps    []OwnedTemp `json:"ownedTemps"`
}

func EmptySnapshot() Snapshot { return Snapshot{SchemaVersion: SchemaVersion} }

func (snapshot Snapshot) Validate() error {
	if snapshot.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: unsupported registry schema", ErrInvalid)
	}
	identities := map[ID]string{}
	serverIDs := map[ID]bool{}
	ownerIDs := map[ID]bool{}
	tombstonedTargets := map[ID]bool{}
	claim := func(id ID, kind string) error {
		if previous, exists := identities[id]; exists {
			return fmt.Errorf("%w: %s and %s use %s", ErrDuplicate, previous, kind, id)
		}
		identities[id] = kind
		return nil
	}
	for _, server := range snapshot.Servers {
		if err := errors.Join(server.Validate(), claim(server.ID, "server")); err != nil {
			return err
		}
		serverIDs[server.ID] = true
		ownerIDs[server.ID] = true
	}
	for _, item := range snapshot.Deliveries {
		if err := errors.Join(item.Validate(), claim(item.ID, "delivery")); err != nil {
			return err
		}
		if !serverIDs[item.ServerID] {
			return fmt.Errorf("%w: delivery server %s", ErrNotFound, item.ServerID)
		}
		ownerIDs[item.ID] = true
	}
	for _, tombstone := range snapshot.Tombstones {
		if err := errors.Join(tombstone.Validate(), claim(tombstone.ID, "tombstone")); err != nil {
			return err
		}
		if previous, exists := identities[tombstone.TargetID]; exists {
			return fmt.Errorf("%w: tombstoned target is still registered as %s", ErrDuplicate, previous)
		}
		if tombstonedTargets[tombstone.TargetID] {
			return fmt.Errorf("%w: target %s has multiple tombstones", ErrDuplicate, tombstone.TargetID)
		}
		tombstonedTargets[tombstone.TargetID] = true
		ownerIDs[tombstone.TargetID] = true
	}
	for _, temporary := range snapshot.OwnedTemps {
		if err := errors.Join(temporary.Validate(), claim(temporary.ID, "owned temporary resource")); err != nil {
			return err
		}
		if !ownerIDs[temporary.OwnerID] {
			return fmt.Errorf("%w: temporary owner %s", ErrNotFound, temporary.OwnerID)
		}
	}
	return nil
}

type Registry struct{ snapshot Snapshot }

func NewRegistry(snapshot Snapshot) (*Registry, error) {
	if err := snapshot.Validate(); err != nil {
		return nil, err
	}
	return &Registry{snapshot: cloneSnapshot(snapshot)}, nil
}

func (registry *Registry) Snapshot() Snapshot { return cloneSnapshot(registry.snapshot) }

func (registry *Registry) RegisterServer(server Server) error {
	if err := server.Validate(); err != nil {
		return err
	}
	if registry.hasIdentity(server.ID) {
		return fmt.Errorf("%w: %s", ErrDuplicate, server.ID)
	}
	registry.snapshot.Servers = append(registry.snapshot.Servers, server)
	return nil
}

func (registry *Registry) RegisterDelivery(item Delivery) error {
	if err := item.Validate(); err != nil {
		return err
	}
	if registry.hasIdentity(item.ID) {
		return fmt.Errorf("%w: %s", ErrDuplicate, item.ID)
	}
	if registry.serverIndex(item.ServerID) < 0 {
		return fmt.Errorf("%w: server %s", ErrNotFound, item.ServerID)
	}
	registry.snapshot.Deliveries = append(registry.snapshot.Deliveries, item)
	return nil
}

func (registry *Registry) TransitionServer(id ID, state State, at time.Time) error {
	index := registry.serverIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: server %s", ErrNotFound, id)
	}
	server := &registry.snapshot.Servers[index]
	if !CanTransition(server.State, state) {
		return fmt.Errorf("%w: server transition %s to %s", ErrInvalid, server.State, state)
	}
	if at.IsZero() || at.Before(server.UpdatedAt) {
		return fmt.Errorf("%w: server update time regressed", ErrInvalid)
	}
	server.State = state
	server.UpdatedAt = at
	return nil
}

func (registry *Registry) UpdatePolicy(id ID, expected uint64, policy Policy, at time.Time) error {
	index := registry.deliveryIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: delivery %s", ErrNotFound, id)
	}
	current := registry.snapshot.Deliveries[index].Policy.Version
	if current != expected || policy.Version != expected+1 {
		return fmt.Errorf("%w: policy expected %d, current %d, next %d", ErrRevisionConflict, expected, current, policy.Version)
	}
	if err := policy.Validate(); err != nil {
		return err
	}
	if at.IsZero() || at.Before(registry.snapshot.Deliveries[index].UpdatedAt) {
		return fmt.Errorf("%w: policy update time regressed", ErrInvalid)
	}
	registry.snapshot.Deliveries[index].Policy = clonePolicy(policy)
	registry.snapshot.Deliveries[index].UpdatedAt = at
	return nil
}

func (registry *Registry) TransitionDelivery(id ID, state State, at time.Time) error {
	index := registry.deliveryIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: delivery %s", ErrNotFound, id)
	}
	item := &registry.snapshot.Deliveries[index]
	if !CanTransition(item.State, state) {
		return fmt.Errorf("%w: delivery transition %s to %s", ErrInvalid, item.State, state)
	}
	if at.IsZero() || at.Before(item.UpdatedAt) {
		return fmt.Errorf("%w: delivery update time regressed", ErrInvalid)
	}
	item.State = state
	item.UpdatedAt = at
	return nil
}

func (registry *Registry) AddOwnedTemp(temporary OwnedTemp) error {
	if err := temporary.Validate(); err != nil {
		return err
	}
	if registry.hasIdentity(temporary.ID) {
		return fmt.Errorf("%w: %s", ErrDuplicate, temporary.ID)
	}
	if !registry.hasOwner(temporary.OwnerID) {
		return fmt.Errorf("%w: owner %s", ErrNotFound, temporary.OwnerID)
	}
	registry.snapshot.OwnedTemps = append(registry.snapshot.OwnedTemps, temporary)
	return nil
}

func (registry *Registry) RemoveOwnedTemp(id ID) error {
	for index, temporary := range registry.snapshot.OwnedTemps {
		if temporary.ID == id {
			registry.snapshot.OwnedTemps = append(registry.snapshot.OwnedTemps[:index], registry.snapshot.OwnedTemps[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: temporary resource %s", ErrNotFound, id)
}

func (registry *Registry) TombstoneDelivery(id ID, tombstone Tombstone) error {
	index := registry.deliveryIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: delivery %s", ErrNotFound, id)
	}
	if tombstone.TargetID != id || tombstone.Kind != TargetDelivery {
		return fmt.Errorf("%w: mismatched delivery tombstone", ErrInvalid)
	}
	if err := tombstone.Validate(); err != nil {
		return err
	}
	if registry.hasIdentity(tombstone.ID) {
		return fmt.Errorf("%w: %s", ErrDuplicate, tombstone.ID)
	}
	registry.snapshot.Deliveries = append(registry.snapshot.Deliveries[:index], registry.snapshot.Deliveries[index+1:]...)
	registry.snapshot.Tombstones = append(registry.snapshot.Tombstones, tombstone)
	return nil
}

func (registry *Registry) TombstoneServer(id ID, tombstone Tombstone) error {
	index := registry.serverIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: server %s", ErrNotFound, id)
	}
	for _, item := range registry.snapshot.Deliveries {
		if item.ServerID == id {
			return fmt.Errorf("%w: server still owns deliveries", ErrInvalid)
		}
	}
	if tombstone.TargetID != id || tombstone.Kind != TargetServer {
		return fmt.Errorf("%w: mismatched server tombstone", ErrInvalid)
	}
	if err := tombstone.Validate(); err != nil {
		return err
	}
	if registry.hasIdentity(tombstone.ID) {
		return fmt.Errorf("%w: %s", ErrDuplicate, tombstone.ID)
	}
	registry.snapshot.Servers = append(registry.snapshot.Servers[:index], registry.snapshot.Servers[index+1:]...)
	registry.snapshot.Tombstones = append(registry.snapshot.Tombstones, tombstone)
	return nil
}

func (registry *Registry) SetCounters(id ID, counters CounterSnapshot, at time.Time) error {
	index := registry.deliveryIndex(id)
	if index < 0 {
		return fmt.Errorf("%w: delivery %s", ErrNotFound, id)
	}
	current := registry.snapshot.Deliveries[index].Counters
	if err := counters.Validate(); err != nil || counters.Read < current.Read || counters.Sent < current.Sent || counters.Confirmed < current.Confirmed || at.Before(registry.snapshot.Deliveries[index].UpdatedAt) {
		return errors.Join(err, fmt.Errorf("%w: counters or update time regressed", ErrInvalid))
	}
	registry.snapshot.Deliveries[index].Counters = counters
	registry.snapshot.Deliveries[index].UpdatedAt = at
	return nil
}

func (registry *Registry) hasIdentity(id ID) bool {
	for _, server := range registry.snapshot.Servers {
		if server.ID == id {
			return true
		}
	}
	for _, item := range registry.snapshot.Deliveries {
		if item.ID == id {
			return true
		}
	}
	for _, tombstone := range registry.snapshot.Tombstones {
		if tombstone.ID == id || tombstone.TargetID == id {
			return true
		}
	}
	for _, temporary := range registry.snapshot.OwnedTemps {
		if temporary.ID == id {
			return true
		}
	}
	return false
}

func (registry *Registry) hasOwner(id ID) bool {
	if registry.serverIndex(id) >= 0 || registry.deliveryIndex(id) >= 0 {
		return true
	}
	for _, tombstone := range registry.snapshot.Tombstones {
		if tombstone.TargetID == id {
			return true
		}
	}
	return false
}

func (registry *Registry) serverIndex(id ID) int {
	for index, server := range registry.snapshot.Servers {
		if server.ID == id {
			return index
		}
	}
	return -1
}

func (registry *Registry) deliveryIndex(id ID) int {
	for index, item := range registry.snapshot.Deliveries {
		if item.ID == id {
			return index
		}
	}
	return -1
}

func cloneSnapshot(snapshot Snapshot) Snapshot {
	result := snapshot
	result.Servers = append([]Server(nil), snapshot.Servers...)
	result.Deliveries = append([]Delivery(nil), snapshot.Deliveries...)
	for index := range result.Deliveries {
		result.Deliveries[index].Policy = clonePolicy(result.Deliveries[index].Policy)
	}
	result.Tombstones = append([]Tombstone(nil), snapshot.Tombstones...)
	result.OwnedTemps = append([]OwnedTemp(nil), snapshot.OwnedTemps...)
	sort.Slice(result.Servers, func(i, j int) bool { return result.Servers[i].ID < result.Servers[j].ID })
	sort.Slice(result.Deliveries, func(i, j int) bool { return result.Deliveries[i].ID < result.Deliveries[j].ID })
	sort.Slice(result.Tombstones, func(i, j int) bool { return result.Tombstones[i].ID < result.Tombstones[j].ID })
	sort.Slice(result.OwnedTemps, func(i, j int) bool { return result.OwnedTemps[i].ID < result.OwnedTemps[j].ID })
	return result
}

func clonePolicy(policy Policy) Policy {
	policy.AllowIP = append([]string(nil), policy.AllowIP...)
	sort.Strings(policy.AllowIP)
	return policy
}
