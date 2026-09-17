# admin-ui Specification

## Purpose
Define the isolated, local administration runtime, guarded control API, and accessible browser interface used to inspect and manage live Courier data deliveries without exposing credentials or weakening worker lifecycle boundaries.

## Requirements

### Requirement: Independent administration singleton

Courier SHALL run at most one administration process per user, SHALL control it only through UUID-verified private IPC, and SHALL keep its lifecycle independent from all data workers.

#### Scenario: Data servers are stopped

- **WHEN** an operator runs `courier servers stop --all` while the administration UI is live
- **THEN** Courier stops the data workers without stopping or registering the administration process

### Requirement: Foreground and background administration

Courier SHALL support foreground and explicit background UI startup, SHALL report the canonical local URL only after readiness is verified, and SHALL make `courier ui stop` idempotent without signaling a recorded PID.

#### Scenario: Concurrent starts race

- **WHEN** two commands attempt to start the administration UI concurrently
- **THEN** one singleton owns the private lock and both callers either identify that live instance or receive a deterministic conflict without creating two listeners

### Requirement: Guarded local API

Courier SHALL bind the administration API only to loopback, require the listener's exact Host on every request, require an exact same-origin Origin on mutations, emit no permissive CORS policy, and bound request bodies, timeouts, and event subscribers.

#### Scenario: Browser request has a forged origin

- **WHEN** a policy or stop request uses an absent or foreign Origin
- **THEN** Courier rejects it before contacting a data worker or changing state

### Requirement: Secret-free authoritative administration

Courier SHALL derive administration snapshots and actions from private registry discovery plus authoritative data-worker IPC and SHALL never return credentials, sessions, opaque resource tokens, runtime definitions, or unsafe diagnostics.

#### Scenario: Protected delivery is inspected

- **WHEN** the administration UI loads a password- or Basic-protected delivery
- **THEN** it receives the controlled authentication mode and policy but no credential or protected resource identifier

### Requirement: Optimistic policy control

Courier SHALL apply a policy update only when its expected version matches the live delivery policy and SHALL expose version conflicts distinctly so clients refresh rather than overwrite concurrent changes.

#### Scenario: Two editors update one delivery

- **WHEN** the second editor submits an obsolete expected policy version
- **THEN** Courier returns a conflict and preserves the first committed policy

### Requirement: Live administration events

Courier SHALL provide bounded same-origin server-sent events containing secret-free authoritative snapshots and SHALL release every subscription after disconnect or shutdown.

#### Scenario: Event client disconnects

- **WHEN** a browser closes its event stream
- **THEN** Courier cancels the subscription and retains no client or worker progress resource

### Requirement: Shared accessible administration application

Courier SHALL embed a deterministic React operations application built from the repository-owned shadcn component system with typed English/Russian catalogs, cyclic system/light/dark preference, overview counters, a keyboard-operable server/delivery navigator, a selected route and policy inspector, and explicit refresh, save, and stop actions backed only by the guarded API. Related counters SHALL use one continuous strip and navigator/inspector SHALL use one shared workspace rather than independent bordered cards. Selection SHALL persist by UUID across authoritative SSE snapshots when possible and SHALL fall back deterministically when the selected object disappears. The application SHALL provide no terminal prompt or arbitrary execution surface.

#### Scenario: A snapshot updates the selected delivery

- **WHEN** an SSE snapshot still contains the selected UUID
- **THEN** counters, route, state, and policy update in place while selection and focused controls remain stable

#### Scenario: A selected object disappears

- **WHEN** an authoritative snapshot no longer contains the selected server or delivery
- **THEN** the navigator selects the first valid remaining object or renders the localized empty state without exposing stale data

#### Scenario: The viewport is narrow

- **WHEN** the administration application renders on a mobile viewport
- **THEN** navigator and inspector stack without hiding actions or creating horizontal overflow

#### Scenario: An operator manages a delivery

- **WHEN** the operator selects a delivery
- **THEN** lightweight shadcn navigation and policy controls present current route, counters, policy, save, and stop actions with complete keyboard operation and without nested card-on-card framing

#### Scenario: An operator inspects and changes a delivery

- **WHEN** a live snapshot arrives or a semantic policy, refresh, or stop control is activated
- **THEN** the React operations application updates authoritative secret-free state through the existing guarded API without recreating a shell interface

#### Scenario: Locale and theme change

- **WHEN** an operator cycles to Russian and dark theme
- **THEN** administration content localizes, the shared theme resolves to dark, and both preferences persist through the shared React preference provider
