// Package contract validates Courier's machine-readable public CLI contract.
package contract

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var semanticVersion = regexp.MustCompile(`^[0-9]+[.][0-9]+[.][0-9]+$`)
var contractIdentifier = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Contract is the canonical inventory used by documentation and parity checks.
type Contract struct {
	SchemaVersion   int        `yaml:"schema_version"`
	ContractVersion string     `yaml:"contract_version"`
	TargetRelease   string     `yaml:"target_release"`
	Language        string     `yaml:"language"`
	EndpointKinds   []Endpoint `yaml:"endpoint_kinds"`
	Commands        []Command  `yaml:"commands"`
	Flags           []Flag     `yaml:"flags"`
	Routes          []Route    `yaml:"routes"`
	Unsupported     []string   `yaml:"unsupported"`
	Examples        []string   `yaml:"examples"`
}

type Endpoint struct {
	Name   string `yaml:"name"`
	Status string `yaml:"status"`
	Syntax string `yaml:"syntax"`
}

type Command struct {
	Name      string     `yaml:"name"`
	Path      string     `yaml:"path"`
	Usage     string     `yaml:"usage"`
	Status    string     `yaml:"status"`
	System    bool       `yaml:"system"`
	Arguments []Argument `yaml:"arguments"`
	Flags     []string   `yaml:"flags"`
}

type Argument struct {
	Name         string `yaml:"name"`
	Kind         string `yaml:"kind"`
	Required     bool   `yaml:"required"`
	Prefix       string `yaml:"prefix"`
	OmitWhenFlag string `yaml:"omit_when_flag"`
}

type Flag struct {
	Name        string   `yaml:"name"`
	Syntax      string   `yaml:"syntax"`
	Status      string   `yaml:"status"`
	ValueKind   string   `yaml:"value_kind"`
	Choices     []string `yaml:"choices"`
	Placeholder string   `yaml:"placeholder"`
	Repeatable  bool     `yaml:"repeatable"`
	Default     string   `yaml:"default"`
	AppliesTo   []string `yaml:"applies_to"`
	Conflicts   []string `yaml:"conflicts"`
	Requires    []string `yaml:"requires"`
}

type Route struct {
	Name         string   `yaml:"name"`
	Status       string   `yaml:"status"`
	Source       []string `yaml:"source"`
	Destination  []string `yaml:"destination"`
	AllowedFlags []string `yaml:"allowed_flags"`
}

// Load reads and strictly decodes a contract.
func Load(name string) (Contract, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return Contract{}, err
	}
	var result Contract
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&result); err != nil {
		return Contract{}, err
	}
	if err := result.Validate(); err != nil {
		return Contract{}, err
	}
	return result, nil
}

// Validate checks cross-references and compatibility metadata.
func (c Contract) Validate() error {
	if c.SchemaVersion != 2 || !semanticVersion.MatchString(c.ContractVersion) || !semanticVersion.MatchString(c.TargetRelease) || c.Language != "en" {
		return errors.New("contract metadata is invalid")
	}
	endpoints, err := namedStatuses("endpoint", c.EndpointKinds, func(value Endpoint) (string, string) { return value.Name, value.Status })
	if err != nil {
		return err
	}
	flags, err := namedStatuses("flag", c.Flags, func(value Flag) (string, string) { return value.Name, value.Status })
	if err != nil {
		return err
	}
	routes, err := namedStatuses("route", c.Routes, func(value Route) (string, string) { return value.Name, value.Status })
	if err != nil {
		return err
	}
	commands, err := namedStatuses("command", c.Commands, func(value Command) (string, string) { return value.Name, value.Status })
	if err != nil {
		return err
	}
	paths := map[string]struct{}{}
	for _, command := range c.Commands {
		if command.Path == "" || command.Usage == "" {
			return fmt.Errorf("command %q has incomplete syntax", command.Name)
		}
		if _, exists := paths[command.Path]; exists {
			return fmt.Errorf("duplicate command path %q", command.Path)
		}
		paths[command.Path] = struct{}{}
		if command.System != (command.Status == "system") {
			return fmt.Errorf("command %q has inconsistent system status", command.Name)
		}
		if err := references("command flag", command.Flags, flags); err != nil {
			return fmt.Errorf("command %q: %w", command.Name, err)
		}
		argumentNames := map[string]struct{}{}
		for _, argument := range command.Arguments {
			if !contractIdentifier.MatchString(argument.Name) || !validArgumentKind(argument.Kind) {
				return fmt.Errorf("command %q has invalid argument metadata", command.Name)
			}
			if _, exists := argumentNames[argument.Name]; exists {
				return fmt.Errorf("command %q has duplicate argument %q", command.Name, argument.Name)
			}
			argumentNames[argument.Name] = struct{}{}
			if argument.Prefix != "" && (strings.ContainsAny(argument.Prefix, " \t\r\n") || strings.HasPrefix(argument.Prefix, "-")) {
				return fmt.Errorf("command %q argument %q has invalid prefix", command.Name, argument.Name)
			}
			if argument.OmitWhenFlag != "" && !contains(command.Flags, argument.OmitWhenFlag) {
				return fmt.Errorf("command %q argument %q has unknown suppressing flag %q", command.Name, argument.Name, argument.OmitWhenFlag)
			}
		}
	}
	for _, flag := range c.Flags {
		if flag.Syntax == "" || flag.Default == "" || len(flag.AppliesTo) == 0 || !validValueKind(flag.ValueKind) {
			return fmt.Errorf("flag %q has incomplete behavior", flag.Name)
		}
		if flag.ValueKind == "boolean" {
			if flag.Placeholder != "" || len(flag.Choices) != 0 || flag.Repeatable || flag.Default != "false" {
				return fmt.Errorf("flag %q has invalid boolean behavior", flag.Name)
			}
		} else if flag.Placeholder == "" {
			return fmt.Errorf("flag %q has no value placeholder", flag.Name)
		}
		if flag.ValueKind == "enum" {
			if len(flag.Choices) == 0 || !contains(flag.Choices, flag.Default) {
				return fmt.Errorf("flag %q has invalid enum choices", flag.Name)
			}
			if err := uniqueValues("choice", flag.Choices); err != nil {
				return fmt.Errorf("flag %q: %w", flag.Name, err)
			}
		} else if len(flag.Choices) != 0 {
			return fmt.Errorf("flag %q has choices for non-enum value", flag.Name)
		}
		if err := references("conflict", flag.Conflicts, flags); err != nil {
			return fmt.Errorf("flag %q: %w", flag.Name, err)
		}
		if contains(flag.Conflicts, flag.Name) {
			return fmt.Errorf("flag %q conflicts with itself", flag.Name)
		}
		if err := references("requirement", flag.Requires, flags); err != nil {
			return fmt.Errorf("flag %q: %w", flag.Name, err)
		}
		if contains(flag.Requires, flag.Name) {
			return fmt.Errorf("flag %q requires itself", flag.Name)
		}
	}
	applicability := map[string]struct{}{
		"local-to-local": {}, "local-to-ssh": {}, "ssh-to-local": {}, "ssh-to-ssh": {},
	}
	for name := range routes {
		applicability[name] = struct{}{}
	}
	for name := range commands {
		applicability[name] = struct{}{}
	}
	flagValues := make(map[string]Flag, len(c.Flags))
	for _, flag := range c.Flags {
		if err := references("applicability", flag.AppliesTo, applicability); err != nil {
			return fmt.Errorf("flag %q: %w", flag.Name, err)
		}
		flagValues[flag.Name] = flag
	}
	for _, route := range c.Routes {
		if err := references("source endpoint", route.Source, endpoints); err != nil {
			return fmt.Errorf("route %q: %w", route.Name, err)
		}
		if err := references("destination endpoint", route.Destination, endpoints); err != nil {
			return fmt.Errorf("route %q: %w", route.Name, err)
		}
		if err := references("allowed flag", route.AllowedFlags, flags); err != nil {
			return fmt.Errorf("route %q: %w", route.Name, err)
		}
		for _, name := range route.AllowedFlags {
			if !appliesToRoute(flagValues[name], route.Name) {
				return fmt.Errorf("route %q allows inapplicable flag %q", route.Name, name)
			}
		}
	}
	allowedByRoute := make(map[string]map[string]bool, len(c.Routes))
	for _, route := range c.Routes {
		allowedByRoute[route.Name] = make(map[string]bool, len(route.AllowedFlags))
		for _, name := range route.AllowedFlags {
			allowedByRoute[route.Name][name] = true
		}
	}
	for _, flag := range c.Flags {
		for _, scope := range flag.AppliesTo {
			route := scope
			if isPathDirection(scope) {
				route = "path-to-path"
			}
			if _, isRoute := routes[route]; isRoute && !allowedByRoute[route][flag.Name] {
				return fmt.Errorf("flag %q applies to route %q but is not allowed by it", flag.Name, route)
			}
		}
	}
	if len(commands) == 0 || len(routes) == 0 || len(c.Unsupported) == 0 || len(c.Examples) == 0 {
		return errors.New("contract inventory is incomplete")
	}
	return nil
}

func appliesToRoute(flag Flag, route string) bool {
	for _, scope := range flag.AppliesTo {
		if scope == route || route == "path-to-path" && isPathDirection(scope) {
			return true
		}
	}
	return false
}

func isPathDirection(value string) bool {
	return value == "local-to-local" || value == "local-to-ssh" || value == "ssh-to-local" || value == "ssh-to-ssh"
}

type contractItem interface {
	Endpoint | Flag | Route | Command
}

func namedStatuses[T contractItem](kind string, values []T, fields func(T) (string, string)) (map[string]struct{}, error) {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		name, status := fields(value)
		if name == "" || !validStatus(status) {
			return nil, fmt.Errorf("%s metadata is invalid", kind)
		}
		if _, exists := result[name]; exists {
			return nil, fmt.Errorf("duplicate %s %q", kind, name)
		}
		result[name] = struct{}{}
	}
	return result, nil
}

func validStatus(value string) bool {
	return value == "shipped" || value == "planned" || value == "system"
}

func validArgumentKind(value string) bool {
	return value == "endpoint" || value == "uuid" || value == "command-path"
}

func validValueKind(value string) bool {
	switch value {
	case "boolean", "text", "unsigned", "enum", "host-port", "ip-cidr", "path", "regex", "quantity", "rate":
		return true
	default:
		return false
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func uniqueValues(kind string, values []string) error {
	seen := map[string]struct{}{}
	for _, value := range values {
		if value == "" {
			return fmt.Errorf("empty %s", kind)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("duplicate %s %q", kind, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func references(kind string, values []string, known map[string]struct{}) error {
	seen := map[string]struct{}{}
	for _, value := range values {
		if _, ok := known[value]; !ok {
			return fmt.Errorf("unknown %s %q", kind, value)
		}
		if _, ok := seen[value]; ok {
			return fmt.Errorf("duplicate %s %q", kind, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}
