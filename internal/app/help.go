package app

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type structuredArgumentHelp struct {
	Syntax      string
	Required    bool
	Description string
}

type structuredFlagHelp struct {
	Name        string
	Syntax      string
	Required    bool
	Description string
	Default     string
	Repeatable  bool
	AppliesTo   string
	Requires    []string
	Conflicts   []string
}

type structuredCommandHelp struct {
	Usage     string
	Arguments []structuredArgumentHelp
	Flags     []structuredFlagHelp
}

func installContractHelp(root *cobra.Command) {
	var visit func(*cobra.Command, string)
	visit = func(parent *cobra.Command, prefix string) {
		for _, command := range parent.Commands() {
			path := strings.TrimSpace(prefix + " " + command.Name())
			if metadata, ok := generatedCommandHelp[path]; ok {
				metadata := metadata
				for _, flag := range metadata.Flags {
					if value := command.Flags().Lookup(flag.Name); value != nil {
						value.Usage = flag.Description
					}
				}
				command.SetHelpFunc(func(command *cobra.Command, _ []string) { renderStructuredHelp(command, metadata) })
			}
			visit(command, path)
		}
	}
	visit(root, "")
}

func renderStructuredHelp(command *cobra.Command, metadata structuredCommandHelp) {
	if command.Short != "" {
		fmt.Fprintln(command.OutOrStdout(), command.Short)
		fmt.Fprintln(command.OutOrStdout())
	}
	fmt.Fprintln(command.OutOrStdout(), "Usage:")
	fmt.Fprintf(command.OutOrStdout(), "  %s\n", metadata.Usage)
	if len(metadata.Arguments) != 0 {
		fmt.Fprintln(command.OutOrStdout(), "\nArguments:")
		for _, argument := range metadata.Arguments {
			fmt.Fprintf(command.OutOrStdout(), "  %-28s %s\n      %s\n", argument.Syntax, requirementLabel(argument.Required), argument.Description)
		}
	}
	if len(metadata.Flags) != 0 {
		fmt.Fprintln(command.OutOrStdout(), "\nOptions:")
		for _, flag := range metadata.Flags {
			fmt.Fprintf(command.OutOrStdout(), "  %-28s %s\n      %s\n", flag.Syntax, requirementLabel(flag.Required), flag.Description)
			fmt.Fprintf(command.OutOrStdout(), "      Applies to: %s; default: %s; repeatable: %s", flag.AppliesTo, flag.Default, yesNo(flag.Repeatable))
			if len(flag.Requires) != 0 {
				fmt.Fprintf(command.OutOrStdout(), "; requires: --%s", strings.Join(flag.Requires, ", --"))
			}
			if len(flag.Conflicts) != 0 {
				fmt.Fprintf(command.OutOrStdout(), "; conflicts: --%s", strings.Join(flag.Conflicts, ", --"))
			}
			fmt.Fprintln(command.OutOrStdout())
		}
	}
}

func requirementLabel(required bool) string {
	if required {
		return "Required"
	}
	return "Optional"
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
