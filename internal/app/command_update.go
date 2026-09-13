package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCommand(run func(context.Context) (update.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Install the latest verified GitHub release",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if run == nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: errors.New("updater is unavailable")}
			}
			result, err := run(command.Context())
			if err != nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: err}
			}
			if result.Current {
				fmt.Fprintf(command.OutOrStdout(), "courier %s is current\n", result.From)
			} else {
				fmt.Fprintf(command.OutOrStdout(), "updated courier from %s to %s\n", result.From, result.To)
				if result.Notes != "" {
					fmt.Fprintln(command.OutOrStdout(), result.Notes)
				}
			}
			return nil
		},
	}
}
