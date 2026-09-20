package contract

import "encoding/json"

// LandingData is the generated, public, shipped-only contract projection used
// by the static project site.
type LandingData struct {
	ContractVersion string            `json:"contractVersion"`
	TargetRelease   string            `json:"targetRelease"`
	Endpoints       []LandingEndpoint `json:"endpoints"`
	Commands        []LandingCommand  `json:"commands"`
	Flags           []LandingFlag     `json:"flags"`
	Routes          []LandingRoute    `json:"routes"`
	Examples        []string          `json:"examples"`
}

type LandingEndpoint struct {
	Name   string `json:"name"`
	Syntax string `json:"syntax"`
}

type LandingCommand struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Usage     string            `json:"usage"`
	System    bool              `json:"system"`
	Arguments []LandingArgument `json:"arguments"`
	Flags     []string          `json:"flags"`
}

type LandingArgument struct {
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Required     bool   `json:"required"`
	Prefix       string `json:"prefix,omitempty"`
	OmitWhenFlag string `json:"omitWhenFlag,omitempty"`
}

type LandingFlag struct {
	Name        string   `json:"name"`
	Syntax      string   `json:"syntax"`
	ValueKind   string   `json:"valueKind"`
	Choices     []string `json:"choices"`
	Placeholder string   `json:"placeholder"`
	Repeatable  bool     `json:"repeatable"`
	Default     string   `json:"default"`
	AppliesTo   []string `json:"appliesTo"`
	Conflicts   []string `json:"conflicts"`
	Requires    []string `json:"requires"`
}

type LandingRoute struct {
	Name         string   `json:"name"`
	Source       []string `json:"source"`
	Destination  []string `json:"destination"`
	AllowedFlags []string `json:"allowedFlags"`
}

// Landing projects only shipped and system inventory from the canonical
// contract. Planned entries cannot reach the public site bundle.
func (c Contract) Landing() LandingData {
	result := LandingData{ContractVersion: c.ContractVersion, TargetRelease: c.TargetRelease, Examples: append([]string(nil), c.Examples...)}
	for _, endpoint := range c.EndpointKinds {
		if endpoint.Status == "shipped" {
			result.Endpoints = append(result.Endpoints, LandingEndpoint{Name: endpoint.Name, Syntax: endpoint.Syntax})
		}
	}
	for _, command := range c.Commands {
		if command.Status == "shipped" || command.Status == "system" {
			arguments := make([]LandingArgument, 0, len(command.Arguments))
			for _, argument := range command.Arguments {
				arguments = append(arguments, LandingArgument{Name: argument.Name, Kind: argument.Kind, Required: argument.Required, Prefix: argument.Prefix, OmitWhenFlag: argument.OmitWhenFlag})
			}
			result.Commands = append(result.Commands, LandingCommand{Name: command.Name, Path: command.Path, Usage: command.Usage, System: command.System, Arguments: arguments, Flags: append([]string(nil), command.Flags...)})
		}
	}
	for _, flag := range c.Flags {
		if flag.Status == "shipped" {
			result.Flags = append(result.Flags, LandingFlag{
				Name: flag.Name, Syntax: flag.Syntax, ValueKind: flag.ValueKind, Choices: append([]string(nil), flag.Choices...), Placeholder: flag.Placeholder,
				Repeatable: flag.Repeatable, Default: flag.Default, AppliesTo: append([]string(nil), flag.AppliesTo...),
				Conflicts: append([]string(nil), flag.Conflicts...), Requires: append([]string(nil), flag.Requires...),
			})
		}
	}
	for _, route := range c.Routes {
		if route.Status == "shipped" {
			result.Routes = append(result.Routes, LandingRoute{
				Name: route.Name, Source: append([]string(nil), route.Source...), Destination: append([]string(nil), route.Destination...), AllowedFlags: append([]string(nil), route.AllowedFlags...),
			})
		}
	}
	return result
}

// LandingJSON renders stable, indented JSON with a trailing newline.
func (c Contract) LandingJSON() []byte {
	data, _ := json.MarshalIndent(c.Landing(), "", "  ")
	return append(data, '\n')
}
