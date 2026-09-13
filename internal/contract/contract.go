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
	Name   string   `yaml:"name"`
	Path   string   `yaml:"path"`
	Usage  string   `yaml:"usage"`
	Status string   `yaml:"status"`
	System bool     `yaml:"system"`
	Flags  []string `yaml:"flags"`
}

type Flag struct {
	Name       string   `yaml:"name"`
	Syntax     string   `yaml:"syntax"`
	Status     string   `yaml:"status"`
	Repeatable bool     `yaml:"repeatable"`
	Default    string   `yaml:"default"`
	AppliesTo  []string `yaml:"applies_to"`
	Conflicts  []string `yaml:"conflicts"`
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
	if c.SchemaVersion != 1 || !semanticVersion.MatchString(c.ContractVersion) || !semanticVersion.MatchString(c.TargetRelease) || c.Language != "en" {
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
	}
	for _, flag := range c.Flags {
		if flag.Syntax == "" || flag.Default == "" || len(flag.AppliesTo) == 0 {
			return fmt.Errorf("flag %q has incomplete behavior", flag.Name)
		}
		if err := references("conflict", flag.Conflicts, flags); err != nil {
			return fmt.Errorf("flag %q: %w", flag.Name, err)
		}
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
	}
	if len(commands) == 0 || len(routes) == 0 || len(c.Unsupported) == 0 || len(c.Examples) == 0 {
		return errors.New("contract inventory is incomplete")
	}
	return nil
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
