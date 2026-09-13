package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(identity BuildIdentity) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build date",
		Args:  cobra.NoArgs,
		Run: func(command *cobra.Command, _ []string) {
			fmt.Fprintf(command.OutOrStdout(), "courier %s (commit %s, built %s)\n", identity.Version, identity.Commit, identity.Date)
		},
	}
}
