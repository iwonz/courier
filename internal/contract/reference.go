package contract

import (
	"bytes"
	"fmt"
	"strings"
)

// Reference renders the deterministic human-facing contract reference.
func (c Contract) Reference() []byte {
	var output bytes.Buffer
	fmt.Fprintf(&output, "# Courier CLI reference\n\nContract version `%s`; target release `%s`. This file is generated from `docs/cli-contract.yaml`.\n\n", c.ContractVersion, c.TargetRelease)
	output.WriteString("## Commands\n\n| Status | Command | Kind |\n|---|---|---|\n")
	for _, command := range c.Commands {
		kind := "product"
		if command.System {
			kind = "system"
		}
		fmt.Fprintf(&output, "| %s | `%s` | %s |\n", command.Status, command.Usage, kind)
	}
	output.WriteString("\n## Endpoint kinds\n\n| Status | Kind | Syntax |\n|---|---|---|\n")
	for _, endpoint := range c.EndpointKinds {
		fmt.Fprintf(&output, "| %s | `%s` | %s |\n", endpoint.Status, endpoint.Name, endpoint.Syntax)
	}
	output.WriteString("\n## Routes and options\n\n| Status | Route | Source | Destination | Allowed options |\n|---|---|---|---|---|\n")
	for _, route := range c.Routes {
		flags := "none"
		if len(route.AllowedFlags) != 0 {
			formatted := make([]string, len(route.AllowedFlags))
			for index, flag := range route.AllowedFlags {
				formatted[index] = "`--" + flag + "`"
			}
			flags = strings.Join(formatted, ", ")
		}
		fmt.Fprintf(&output, "| %s | `%s` | %s | %s | %s |\n", route.Status, route.Name, strings.Join(route.Source, ", "), strings.Join(route.Destination, ", "), flags)
	}
	output.WriteString("\n## Options\n\n| Status | Option | Repeatable | Default | Applies to | Conflicts |\n|---|---|---:|---|---|---|\n")
	for _, flag := range c.Flags {
		conflicts := "none"
		if len(flag.Conflicts) != 0 {
			conflicts = strings.Join(flag.Conflicts, ", ")
		}
		fmt.Fprintf(&output, "| %s | `%s` | %t | `%s` | %s | %s |\n", flag.Status, flag.Syntax, flag.Repeatable, flag.Default, strings.Join(flag.AppliesTo, ", "), conflicts)
	}
	output.WriteString("\n## Explicitly unsupported\n\n")
	for _, value := range c.Unsupported {
		fmt.Fprintf(&output, "- `%s`\n", value)
	}
	output.WriteString("\n## Examples\n\n```sh\n")
	for _, example := range c.Examples {
		fmt.Fprintln(&output, example)
	}
	output.WriteString("```\n")
	return output.Bytes()
}
