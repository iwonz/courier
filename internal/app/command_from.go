package app

import (
	"errors"
	"strconv"

	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/selection"
	"github.com/spf13/cobra"
)

func newTransferCommand(dependencies Dependencies) *cobra.Command {
	var archiveMode operation.BoolValue
	var selectionValues operation.OrderedValues
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
			options := make([]operation.Option, 0, 1+len(selectionValues.Options()))
			if value, explicit := archiveMode.Value(); explicit {
				options = append(options, operation.Option{Name: operation.OptionArchive, Value: strconv.FormatBool(value)})
			}
			options = append(options, selectionValues.Options()...)
			plan, err := operation.Build(operation.Request{Source: args[0], Destination: args[2], Options: options})
			if err != nil {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: err}
			}
			if plan.Route != operation.RoutePathToPath {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("requested route is planned but not shipped")}
			}
			selector := selection.Selector(selection.All())
			if len(plan.Options.Selection) != 0 {
				if dependencies.Select == nil {
					return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("selection dependencies are incomplete")}
				}
				selector, err = dependencies.Select(plan.Options.Selection)
				if err != nil {
					return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: err}
				}
			}
			return runTransfer(command.Context(), dependencies, plan, selector, command.OutOrStdout(), command.ErrOrStderr())
		},
	}
	command.Flags().Var(&archiveMode, "archive", "create and transfer <source-name>.tar.gz")
	command.Flags().Lookup("archive").NoOptDefVal = "true"
	command.Flags().Var(selectionValues.For(operation.OptionExclude), "exclude", "exclude a gitignore pattern")
	command.Flags().Var(selectionValues.For(operation.OptionExcludeRegex), "exclude-regex", "exclude paths matching a Go regular expression")
	command.Flags().Var(selectionValues.For(operation.OptionExcludeFrom), "exclude-from", "read gitignore patterns from a local file")
	return command
}
