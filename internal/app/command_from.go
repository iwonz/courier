package app

import (
	"errors"
	"strconv"

	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/spf13/cobra"
)

func newTransferCommand(dependencies Dependencies) *cobra.Command {
	var archiveMode operation.BoolValue
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
			options := make([]operation.Option, 0, 1)
			if value, explicit := archiveMode.Value(); explicit {
				options = append(options, operation.Option{Name: operation.OptionArchive, Value: strconv.FormatBool(value)})
			}
			plan, err := operation.Build(operation.Request{Source: args[0], Destination: args[2], Options: options})
			if err != nil {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: err}
			}
			if plan.Route != operation.RoutePathToPath {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("requested route is planned but not shipped")}
			}
			return runTransfer(command.Context(), dependencies, plan, command.OutOrStdout(), command.ErrOrStderr())
		},
	}
	command.Flags().Var(&archiveMode, "archive", "create and transfer <source-name>.tar.gz")
	command.Flags().Lookup("archive").NoOptDefVal = "true"
	return command
}
