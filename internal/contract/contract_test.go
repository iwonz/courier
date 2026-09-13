package contract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/operation"
	"github.com/spf13/cobra"
)

func validContract() Contract {
	return Contract{
		SchemaVersion: 1, ContractVersion: "0.8.0", TargetRelease: "0.2.0", Language: "en",
		EndpointKinds: []Endpoint{{Name: "local", Status: "shipped", Syntax: "path"}},
		Commands:      []Command{{Name: "from", Path: "from", Usage: "courier from", Status: "shipped", Flags: []string{"archive"}}},
		Flags:         []Flag{{Name: "archive", Syntax: "--archive", Status: "shipped", Default: "false", AppliesTo: []string{"path-to-path"}}},
		Routes:        []Route{{Name: "path-to-path", Status: "shipped", Source: []string{"local"}, Destination: []string{"local"}, AllowedFlags: []string{"archive"}}},
		Unsupported:   []string{"--mirror"}, Examples: []string{"courier from a to b"},
	}
}

func TestLoadAndReference(t *testing.T) {
	directory := t.TempDir()
	name := filepath.Join(directory, "contract.yaml")
	data := "schema_version: 1\ncontract_version: 0.8.0\ntarget_release: 0.2.0\nlanguage: en\nendpoint_kinds:\n  - {name: local, status: shipped, syntax: path}\ncommands:\n  - {name: from, path: from, usage: courier-from, status: shipped, system: false, flags: [archive]}\nflags:\n  - {name: archive, syntax: --archive, status: shipped, repeatable: false, default: 'false', applies_to: [path-to-path], conflicts: []}\nroutes:\n  - {name: path-to-path, status: shipped, source: [local], destination: [local], allowed_flags: [archive]}\nunsupported: [--mirror]\nexamples: [courier-from]\n"
	if err := os.WriteFile(name, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	value, err := Load(name)
	if err != nil {
		t.Fatal(err)
	}
	reference := string(value.Reference())
	for _, text := range []string{"# Courier CLI reference", "| shipped | `courier-from` | product |", "| shipped | `local` | path |", "`--archive`", "Explicitly unsupported", "```sh"} {
		if !strings.Contains(reference, text) {
			t.Fatalf("reference missing %q: %s", text, reference)
		}
	}

	if _, err := Load(filepath.Join(directory, "missing")); err == nil {
		t.Fatal("expected read error")
	}
	if err := os.WriteFile(name, []byte("unknown: value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(name); err == nil {
		t.Fatal("expected strict decode error")
	}
	if err := os.WriteFile(name, []byte("schema_version: 0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(name); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateFailures(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Contract)
	}{
		{"metadata", func(c *Contract) { c.Language = "ru" }},
		{"endpoint metadata", func(c *Contract) { c.EndpointKinds[0].Status = "invalid" }},
		{"duplicate endpoint", func(c *Contract) { c.EndpointKinds = append(c.EndpointKinds, c.EndpointKinds[0]) }},
		{"flag metadata", func(c *Contract) { c.Flags[0].Name = "" }},
		{"duplicate flag", func(c *Contract) { c.Flags = append(c.Flags, c.Flags[0]) }},
		{"route metadata", func(c *Contract) { c.Routes[0].Name = "" }},
		{"duplicate route", func(c *Contract) { c.Routes = append(c.Routes, c.Routes[0]) }},
		{"command metadata", func(c *Contract) { c.Commands[0].Name = "" }},
		{"duplicate command", func(c *Contract) { c.Commands = append(c.Commands, c.Commands[0]) }},
		{"command syntax", func(c *Contract) { c.Commands[0].Usage = "" }},
		{"command path", func(c *Contract) {
			c.Commands = append(c.Commands, Command{Name: "other", Path: "from", Usage: "other", Status: "planned"})
		}},
		{"system status", func(c *Contract) { c.Commands[0].System = true }},
		{"unknown command flag", func(c *Contract) { c.Commands[0].Flags = []string{"missing"} }},
		{"duplicate command flag", func(c *Contract) { c.Commands[0].Flags = []string{"archive", "archive"} }},
		{"flag behavior", func(c *Contract) { c.Flags[0].AppliesTo = nil }},
		{"unknown conflict", func(c *Contract) { c.Flags[0].Conflicts = []string{"missing"} }},
		{"duplicate conflict", func(c *Contract) { c.Flags[0].Conflicts = []string{"archive", "archive"} }},
		{"unknown applicability", func(c *Contract) { c.Flags[0].AppliesTo = []string{"missing"} }},
		{"duplicate applicability", func(c *Contract) { c.Flags[0].AppliesTo = []string{"path-to-path", "path-to-path"} }},
		{"source endpoint", func(c *Contract) { c.Routes[0].Source = []string{"missing"} }},
		{"destination endpoint", func(c *Contract) { c.Routes[0].Destination = []string{"missing"} }},
		{"route flag", func(c *Contract) { c.Routes[0].AllowedFlags = []string{"missing"} }},
		{"inapplicable route flag", func(c *Contract) { c.Flags[0].AppliesTo = []string{"from"} }},
		{"missing route flag", func(c *Contract) { c.Routes[0].AllowedFlags = nil }},
		{"inventory", func(c *Contract) { c.Unsupported = nil }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := validContract()
			test.mutate(&value)
			if err := value.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func plannerContractAndMatrix() (Contract, operation.Matrix) {
	value := validContract()
	matrix := operation.Matrix{
		EndpointKinds: []operation.MatrixEndpoint{{Kind: "local"}},
		Routes:        []operation.RouteDefinition{{Route: operation.RoutePathToPath, Sources: []endpoint.Kind{endpoint.KindLocal}, Destinations: []endpoint.Kind{endpoint.KindLocal}}},
		Options:       map[operation.OptionName][]string{operation.OptionArchive: {"path-to-path"}},
		Allowed:       map[operation.Route][]operation.OptionName{operation.RoutePathToPath: {operation.OptionArchive}},
	}
	return value, matrix
}

func TestCheckPlanner(t *testing.T) {
	value, matrix := plannerContractAndMatrix()
	if err := value.CheckPlanner(matrix); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		mutate func(*Contract, *operation.Matrix)
	}{
		{"endpoint", func(_ *Contract, matrix *operation.Matrix) { matrix.EndpointKinds[0].Kind = "ssh" }},
		{"duplicate route", func(_ *Contract, matrix *operation.Matrix) { matrix.Routes = append(matrix.Routes, matrix.Routes[0]) }},
		{"missing route", func(_ *Contract, matrix *operation.Matrix) { matrix.Routes = nil }},
		{"route source", func(_ *Contract, matrix *operation.Matrix) {
			matrix.Routes[0].Sources = []endpoint.Kind{endpoint.KindSSH}
		}},
		{"route destination", func(_ *Contract, matrix *operation.Matrix) {
			matrix.Routes[0].Destinations = []endpoint.Kind{endpoint.KindSSH}
		}},
		{"route option", func(_ *Contract, matrix *operation.Matrix) { matrix.Allowed[operation.RoutePathToPath] = nil }},
		{"undocumented route", func(contract *Contract, _ *operation.Matrix) { contract.Routes = nil }},
		{"missing option", func(_ *Contract, matrix *operation.Matrix) { delete(matrix.Options, operation.OptionArchive) }},
		{"option scope", func(_ *Contract, matrix *operation.Matrix) {
			matrix.Options[operation.OptionArchive] = []string{"path-to-web"}
		}},
		{"undocumented option", func(contract *Contract, _ *operation.Matrix) { contract.Flags = nil }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, matrix := plannerContractAndMatrix()
			test.mutate(&value, &matrix)
			if err := value.CheckPlanner(matrix); err == nil {
				t.Fatal("expected planner parity error")
			}
		})
	}
}

func TestReferenceEmptyAndSystemValues(t *testing.T) {
	value := validContract()
	value.Commands[0].System = true
	value.Commands[0].Status = "system"
	value.Routes[0].AllowedFlags = nil
	value.Flags[0].Conflicts = []string{"archive"}
	reference := string(value.Reference())
	if !strings.Contains(reference, "| system |") || !strings.Contains(reference, "| none |") || !strings.Contains(reference, "| archive |") {
		t.Fatalf("reference=%s", reference)
	}
}

func TestDirectionApplicability(t *testing.T) {
	value := validContract()
	value.Flags[0].AppliesTo = []string{"local-to-ssh"}
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestCheckCobra(t *testing.T) {
	value := validContract()
	root := &cobra.Command{Use: "courier"}
	from := &cobra.Command{Use: "from", Run: func(*cobra.Command, []string) {}}
	from.Flags().Bool("archive", false, "archive")
	root.AddCommand(from)
	if err := value.CheckCobra(root); err != nil {
		t.Fatal(err)
	}

	missing := validContract()
	missing.Commands[0].Path = "other"
	if err := missing.CheckCobra(root); err == nil {
		t.Fatal("expected command mismatch")
	}

	extra := validContract()
	from.Flags().Bool("extra", false, "extra")
	if err := extra.CheckCobra(root); err == nil {
		t.Fatal("expected flag mismatch")
	}

	broken := validContract()
	broken.Commands[0].Flags = nil
	root = &cobra.Command{Use: "courier"}
	parent := &cobra.Command{Use: "from"}
	parent.AddCommand(&cobra.Command{Use: "child", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(parent)
	if err := broken.CheckCobra(root); err == nil {
		t.Fatal("expected nested command mismatch")
	}
}
