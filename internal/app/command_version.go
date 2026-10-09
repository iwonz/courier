package app

import (
	"fmt"

	"github.com/iwonz/courier/internal/terminalui"
	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
)

func newVersionCommand(identity BuildIdentity) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build date",
		Args:  cobra.NoArgs,
		Run: func(command *cobra.Command, _ []string) {
			releaseURL := update.ReleaseURL("", identity.Version)
			mode := terminalMode(command.OutOrStdout())
			if !mode.Interactive {
				fmt.Fprintf(command.OutOrStdout(), "courier %s (commit %s, built %s)\nrelease: %s\n", identity.Version, identity.Commit, identity.Date, releaseURL)
				return
			}
			panel := terminalui.New(command.OutOrStdout(), mode).Panel("Courier", terminalui.ToneInfo, []terminalui.Field{
				{Label: "Version", Value: identity.Version}, {Label: "Release", Value: releaseURL}, {Label: "Commit", Value: identity.Commit}, {Label: "Built", Value: identity.Date},
			}, "From here to anywhere")
			fmt.Fprint(command.OutOrStdout(), panel)
		},
	}
}
