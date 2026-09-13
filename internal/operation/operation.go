// Package operation builds side-effect-free execution plans from CLI input.
package operation

import (
	"errors"
	"fmt"

	"github.com/iwonz/courier/internal/endpoint"
)

var (
	// ErrInvalidEndpoint identifies endpoint syntax failures.
	ErrInvalidEndpoint = errors.New("invalid operation endpoint")
	// ErrUnsupportedRoute identifies unsupported endpoint directions.
	ErrUnsupportedRoute = errors.New("unsupported operation route")
	// ErrInvalidOption identifies option syntax and matrix failures.
	ErrInvalidOption = errors.New("invalid operation option")
)

// Error provides a stable preflight classification and field name.
type Error struct {
	Class error
	Field string
	Cause error
}

func (e *Error) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("%v: %v", e.Class, e.Cause)
	}
	return fmt.Sprintf("%v %q: %v", e.Class, e.Field, e.Cause)
}

// Unwrap exposes both the stable class and the underlying detail.
func (e *Error) Unwrap() []error { return []error{e.Class, e.Cause} }

// Route identifies the runtime family selected by endpoint kinds.
type Route uint8

const (
	RoutePathToPath Route = iota
	RouteWebToPath
	RoutePathToWeb
	RouteWebhookToPath
	RoutePathToHTTP
)

// String returns the canonical contract route name.
func (r Route) String() string {
	switch r {
	case RoutePathToPath:
		return "path-to-path"
	case RouteWebToPath:
		return "web-to-path"
	case RoutePathToWeb:
		return "path-to-web"
	case RouteWebhookToPath:
		return "webhook-to-path"
	case RoutePathToHTTP:
		return "path-to-http"
	default:
		return "unknown"
	}
}

// Direction identifies the concrete path transport legs in a route.
type Direction uint8

const (
	DirectionNone Direction = iota
	DirectionLocalToLocal
	DirectionLocalToSSH
	DirectionSSHToLocal
	DirectionSSHToSSH
)

// String returns the canonical path-direction name.
func (d Direction) String() string {
	switch d {
	case DirectionNone:
		return "none"
	case DirectionLocalToLocal:
		return "local-to-local"
	case DirectionLocalToSSH:
		return "local-to-ssh"
	case DirectionSSHToLocal:
		return "ssh-to-local"
	case DirectionSSHToSSH:
		return "ssh-to-ssh"
	default:
		return "unknown"
	}
}

// Request is the complete input to preflight planning.
type Request struct {
	Source      string
	Destination string
	Options     []Option
}

// Plan is immutable input for a later route-specific runtime.
type Plan struct {
	Source      endpoint.Endpoint
	Destination endpoint.Endpoint
	Route       Route
	Direction   Direction
	Options     Options
}

// RouteDefinition is one explicit source-kind by destination-kind matrix row.
type RouteDefinition struct {
	Route        Route
	Sources      []endpoint.Kind
	Destinations []endpoint.Kind
}

var routeDefinitions = []RouteDefinition{
	{Route: RoutePathToPath, Sources: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}, Destinations: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}},
	{Route: RouteWebToPath, Sources: []endpoint.Kind{endpoint.KindWeb}, Destinations: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}},
	{Route: RoutePathToWeb, Sources: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}, Destinations: []endpoint.Kind{endpoint.KindWeb}},
	{Route: RouteWebhookToPath, Sources: []endpoint.Kind{endpoint.KindWebhook}, Destinations: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}},
	{Route: RoutePathToHTTP, Sources: []endpoint.Kind{endpoint.KindLocal, endpoint.KindSSH}, Destinations: []endpoint.Kind{endpoint.KindHTTP}},
}

// Build parses and validates a request without performing I/O.
func Build(request Request) (Plan, error) {
	source, err := endpoint.Parse(request.Source)
	if err != nil {
		return Plan{}, &Error{Class: ErrInvalidEndpoint, Field: "source", Cause: err}
	}
	destination, err := endpoint.Parse(request.Destination)
	if err != nil {
		return Plan{}, &Error{Class: ErrInvalidEndpoint, Field: "destination", Cause: err}
	}
	route, direction, err := classify(source, destination)
	if err != nil {
		return Plan{}, err
	}
	options, err := buildOptions(route, direction, request.Options)
	if err != nil {
		return Plan{}, err
	}
	return Plan{Source: source, Destination: destination, Route: route, Direction: direction, Options: options}, nil
}

func classify(source, destination endpoint.Endpoint) (Route, Direction, error) {
	for _, definition := range routeDefinitions {
		if containsKind(definition.Sources, source.Kind) && containsKind(definition.Destinations, destination.Kind) {
			direction := DirectionNone
			if definition.Route == RoutePathToPath {
				direction = pathDirection(source.Kind, destination.Kind)
			}
			return definition.Route, direction, nil
		}
	}
	cause := fmt.Errorf("%s to %s", source.Kind, destination.Kind)
	return 0, 0, &Error{Class: ErrUnsupportedRoute, Cause: cause}
}

func containsKind(values []endpoint.Kind, wanted endpoint.Kind) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func pathDirection(source, destination endpoint.Kind) Direction {
	if source == endpoint.KindSSH {
		if destination == endpoint.KindSSH {
			return DirectionSSHToSSH
		}
		return DirectionSSHToLocal
	}
	if destination == endpoint.KindSSH {
		return DirectionLocalToSSH
	}
	return DirectionLocalToLocal
}
