package app

import (
	"errors"

	"github.com/iwonz/courier/internal/progress"
	"github.com/spf13/cobra"
)

func newTransferCommand(dependencies Dependencies) *cobra.Command {
	archiveMode := false
	command := &cobra.Command{
		Use:   "from <source> to <destination>",
		Short: "Transfer a file or directory",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 3 || args[1] != "to" {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("expected: courier from <source> to <destination> [flags]")}
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			return runTransfer(command.Context(), dependencies, args[0], args[2], archiveMode, command.OutOrStdout(), command.ErrOrStderr())
		},
	}
	command.Flags().BoolVar(&archiveMode, "archive", false, "create and transfer <source-name>.tar.gz")
	return command
}
