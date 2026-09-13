package operation

import (
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iwonz/courier/internal/endpoint"
)

// OptionName is a stable command-contract option identifier.
type OptionName string

const (
	OptionArchive          OptionName = "archive"
	OptionExtract          OptionName = "extract"
	OptionListen           OptionName = "listen"
	OptionBackground       OptionName = "background"
	OptionAuth             OptionName = "auth"
	OptionAuthAttempts     OptionName = "auth-attempts"
	OptionAuthFailAction   OptionName = "auth-fail-action"
	OptionLimit            OptionName = "limit"
	OptionNoUI             OptionName = "no-ui"
	OptionAllowIP          OptionName = "allow-ip"
	OptionExclude          OptionName = "exclude"
	OptionExcludeRegex     OptionName = "exclude-regex"
	OptionExcludeFrom      OptionName = "exclude-from"
	OptionMaxFileSize      OptionName = "max-file-size"
	OptionMaxExtractedSize OptionName = "max-extracted-size"
	OptionUploadRate       OptionName = "upload-rate"
	OptionDownloadRate     OptionName = "download-rate"
	OptionAll              OptionName = "all"
)

// Option is one raw command-line occurrence. Slice order is command-line order.
type Option struct {
	Name  OptionName
	Value string
}

// Occurrence records normalized command-line position for diagnostics and
// ordering-sensitive consumers.
type Occurrence struct {
	Option
	Position int
}

// AuthMode controls HTTP delivery authentication.
type AuthMode string

const (
	AuthNone     AuthMode = "none"
	AuthBasic    AuthMode = "basic"
	AuthPassword AuthMode = "password"
)

// AuthFailAction selects behavior at the authentication-attempt threshold.
type AuthFailAction string

const (
	AuthFailBan  AuthFailAction = "ban"
	AuthFailStop AuthFailAction = "stop"
)

// Limit is either unlimited or a positive delivery count.
type Limit struct {
	Unlimited bool
	Value     uint64
}

// SelectionKind distinguishes ordered selection-rule sources.
type SelectionKind string

const (
	SelectionGitignore SelectionKind = "gitignore"
	SelectionRegex     SelectionKind = "regex"
	SelectionFile      SelectionKind = "file"
)

// SelectionRule is an ordered, validated rule occurrence.
type SelectionRule struct {
	Kind     SelectionKind
	Value    string
	Position int
}

// Options contains effective typed values plus explicit occurrence metadata.
type Options struct {
	Archive          bool
	Extract          bool
	Listen           string
	Background       bool
	Auth             AuthMode
	AuthAttempts     uint64
	AuthFailAction   AuthFailAction
	Limit            Limit
	NoUI             bool
	AllowIP          []netip.Prefix
	Selection        []SelectionRule
	MaxFileSize      Quantity
	MaxExtractedSize Quantity
	UploadRate       Quantity
	DownloadRate     Quantity
	Occurrences      []Occurrence
	explicit         map[OptionName]bool
}

// Explicit reports whether an option occurred on the command line.
func (o Options) Explicit(name OptionName) bool { return o.explicit[name] }

type optionDefinition struct {
	repeatable bool
	routes     map[Route]bool
	directions map[Direction]bool
}

// Matrix is the runtime endpoint, route, and option-applicability inventory
// checked against docs/cli-contract.yaml.
type Matrix struct {
	EndpointKinds []MatrixEndpoint
	Routes        []RouteDefinition
	Options       map[OptionName][]string
	Allowed       map[Route][]OptionName
}

// MatrixEndpoint keeps the public matrix independent from parser
// implementation details while retaining a typed kind.
type MatrixEndpoint struct {
	Kind string
}

var optionDefinitions = map[OptionName]optionDefinition{
	OptionArchive:          {routes: routeSet(RoutePathToPath, RoutePathToWeb, RoutePathToHTTP)},
	OptionExtract:          {routes: routeSet(RoutePathToPath, RouteWebToPath, RouteWebhookToPath)},
	OptionListen:           {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionBackground:       {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionAuth:             {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath, RoutePathToHTTP)},
	OptionAuthAttempts:     {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionAuthFailAction:   {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionLimit:            {routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionNoUI:             {routes: routeSet(RoutePathToWeb)},
	OptionAllowIP:          {repeatable: true, routes: routeSet(RouteWebToPath, RoutePathToWeb, RouteWebhookToPath)},
	OptionExclude:          {repeatable: true, routes: allRoutes()},
	OptionExcludeRegex:     {repeatable: true, routes: allRoutes()},
	OptionExcludeFrom:      {repeatable: true, routes: allRoutes()},
	OptionMaxFileSize:      {routes: routeSet(RouteWebToPath, RouteWebhookToPath)},
	OptionMaxExtractedSize: {routes: routeSet(RoutePathToPath, RouteWebToPath, RouteWebhookToPath)},
	OptionUploadRate:       {routes: routeSet(RouteWebToPath, RouteWebhookToPath, RoutePathToHTTP), directions: directionSet(DirectionLocalToSSH, DirectionSSHToSSH)},
	OptionDownloadRate:     {routes: routeSet(RoutePathToWeb), directions: directionSet(DirectionSSHToLocal, DirectionSSHToSSH)},
	OptionAll:              {},
}

// ContractMatrix returns a defensive snapshot for contract verification.
func ContractMatrix() Matrix {
	result := Matrix{
		EndpointKinds: []MatrixEndpoint{{Kind: "local"}, {Kind: "ssh"}, {Kind: "web"}, {Kind: "webhook"}, {Kind: "http"}},
		Routes:        make([]RouteDefinition, len(routeDefinitions)),
		Options:       make(map[OptionName][]string, len(optionDefinitions)),
		Allowed:       make(map[Route][]OptionName, len(routeDefinitions)),
	}
	for index, definition := range routeDefinitions {
		result.Routes[index] = RouteDefinition{Route: definition.Route, Sources: append([]endpoint.Kind(nil), definition.Sources...), Destinations: append([]endpoint.Kind(nil), definition.Destinations...)}
	}
	for name, definition := range optionDefinitions {
		scopes := make([]string, 0, len(definition.routes)+len(definition.directions))
		for route := range definition.routes {
			scopes = append(scopes, route.String())
		}
		for direction := range definition.directions {
			scopes = append(scopes, direction.String())
		}
		sort.Strings(scopes)
		result.Options[name] = scopes
		for _, route := range routeDefinitions {
			if definition.routes[route.Route] || route.Route == RoutePathToPath && len(definition.directions) != 0 {
				result.Allowed[route.Route] = append(result.Allowed[route.Route], name)
			}
		}
	}
	for route := range result.Allowed {
		sort.Slice(result.Allowed[route], func(left, right int) bool { return result.Allowed[route][left] < result.Allowed[route][right] })
	}
	return result
}

func routeSet(routes ...Route) map[Route]bool {
	result := make(map[Route]bool, len(routes))
	for _, route := range routes {
		result[route] = true
	}
	return result
}

func allRoutes() map[Route]bool {
	return routeSet(RoutePathToPath, RouteWebToPath, RoutePathToWeb, RouteWebhookToPath, RoutePathToHTTP)
}

func directionSet(directions ...Direction) map[Direction]bool {
	result := make(map[Direction]bool, len(directions))
	for _, direction := range directions {
		result[direction] = true
	}
	return result
}

func defaultOptions() Options {
	return Options{
		Listen:           "127.0.0.1:8080",
		Auth:             AuthNone,
		AuthAttempts:     5,
		AuthFailAction:   AuthFailBan,
		Limit:            Limit{Unlimited: true},
		MaxFileSize:      MustParseQuantity("10GiB", false),
		MaxExtractedSize: MustParseQuantity("100GiB", false),
		UploadRate:       Quantity{Unlimited: true, Rate: true},
		DownloadRate:     Quantity{Unlimited: true, Rate: true},
		explicit:         map[OptionName]bool{},
	}
}

func buildOptions(route Route, direction Direction, raw []Option) (Options, error) {
	result := defaultOptions()
	for position, value := range raw {
		definition, known := optionDefinitions[value.Name]
		if !known {
			return Options{}, optionError(value.Name, "unknown option")
		}
		if result.explicit[value.Name] && !definition.repeatable {
			return Options{}, optionError(value.Name, "option is not repeatable")
		}
		if !definition.routes[route] && !definition.directions[direction] {
			return Options{}, optionError(value.Name, fmt.Sprintf("not applicable to %s", route))
		}
		result.explicit[value.Name] = true
		occurrence := Occurrence{Option: value, Position: position}
		result.Occurrences = append(result.Occurrences, occurrence)
		if err := applyOption(&result, occurrence); err != nil {
			return Options{}, err
		}
	}
	if result.Archive && result.Extract {
		return Options{}, optionError(OptionArchive, "conflicts with --extract")
	}
	if result.Explicit(OptionMaxExtractedSize) && !result.Extract {
		return Options{}, optionError(OptionMaxExtractedSize, "requires --extract")
	}
	if result.Auth == AuthPassword && route != RouteWebToPath && route != RoutePathToWeb {
		return Options{}, optionError(OptionAuth, "password mode requires a browser route")
	}
	return result, nil
}

func applyOption(result *Options, occurrence Occurrence) error {
	value := occurrence.Value
	var err error
	switch occurrence.Name {
	case OptionArchive:
		result.Archive, err = parseBool(value)
	case OptionExtract:
		result.Extract, err = parseBool(value)
	case OptionListen:
		err = validateListen(value)
		result.Listen = value
	case OptionBackground:
		result.Background, err = parseBool(value)
	case OptionAuth:
		result.Auth = AuthMode(value)
		if result.Auth != AuthNone && result.Auth != AuthBasic && result.Auth != AuthPassword {
			err = fmt.Errorf("expected none, basic, or password")
		}
	case OptionAuthAttempts:
		result.AuthAttempts, err = positiveInteger(value)
	case OptionAuthFailAction:
		result.AuthFailAction = AuthFailAction(value)
		if result.AuthFailAction != AuthFailBan && result.AuthFailAction != AuthFailStop {
			err = fmt.Errorf("expected ban or stop")
		}
	case OptionLimit:
		result.Limit, err = parseLimit(value)
	case OptionNoUI:
		result.NoUI, err = parseBool(value)
	case OptionAllowIP:
		var prefix netip.Prefix
		prefix, err = parsePrefix(value)
		if err == nil {
			result.AllowIP = append(result.AllowIP, prefix)
		}
	case OptionExclude, OptionExcludeRegex, OptionExcludeFrom:
		err = appendSelection(result, occurrence)
	case OptionMaxFileSize:
		result.MaxFileSize, err = ParseQuantity(value, false)
	case OptionMaxExtractedSize:
		result.MaxExtractedSize, err = ParseQuantity(value, false)
	case OptionUploadRate:
		result.UploadRate, err = ParseQuantity(value, true)
	case OptionDownloadRate:
		result.DownloadRate, err = ParseQuantity(value, true)
	default:
		err = fmt.Errorf("option has no operation value")
	}
	if err != nil {
		return optionError(occurrence.Name, err.Error())
	}
	return nil
}

func parseBool(value string) (bool, error) {
	if value == "" {
		value = "true"
	}
	return strconv.ParseBool(value)
}

func positiveInteger(value string) (uint64, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, fmt.Errorf("expected a positive integer")
	}
	return parsed, nil
}

func parseLimit(value string) (Limit, error) {
	if value == "unlimited" {
		return Limit{Unlimited: true}, nil
	}
	parsed, err := positiveInteger(value)
	return Limit{Value: parsed}, err
}

func validateListen(value string) error {
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("expected host:port")
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return fmt.Errorf("expected port 1..65535")
	}
	return nil
}

func parsePrefix(value string) (netip.Prefix, error) {
	if prefix, err := netip.ParsePrefix(value); err == nil {
		return prefix.Masked(), nil
	}
	address, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("expected IP or CIDR")
	}
	return netip.PrefixFrom(address, address.BitLen()), nil
}

func appendSelection(result *Options, occurrence Occurrence) error {
	if occurrence.Value == "" {
		return fmt.Errorf("value must not be empty")
	}
	kind := SelectionGitignore
	if occurrence.Name == OptionExcludeRegex {
		kind = SelectionRegex
		if _, err := regexp.Compile(occurrence.Value); err != nil {
			return fmt.Errorf("invalid regular expression: %v", err)
		}
	} else if occurrence.Name == OptionExcludeFrom {
		kind = SelectionFile
	}
	result.Selection = append(result.Selection, SelectionRule{Kind: kind, Value: occurrence.Value, Position: occurrence.Position})
	return nil
}

func optionError(name OptionName, message string) error {
	return &Error{Class: ErrInvalidOption, Field: "--" + string(name), Cause: fmt.Errorf("%s", strings.TrimSpace(message))}
}
