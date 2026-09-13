package app

import (
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/selection"
	"github.com/spf13/cobra"
)

type transferFlagValues struct {
	archive, extract, background, noUI                      operation.BoolValue
	listen, auth, authAttempts, authFailAction, limit       operation.SingleValue
	maxFileSize, maxExtractedSize, uploadRate, downloadRate operation.SingleValue
	allowIP                                                 operation.RepeatedValue
	selection                                               operation.OrderedValues
}

func (values *transferFlagValues) options() []operation.Option {
	options := make([]operation.Option, 0, 16+len(values.selection.Options()))
	appendBool := func(name operation.OptionName, value *operation.BoolValue) {
		if parsed, explicit := value.Value(); explicit {
			options = append(options, operation.Option{Name: name, Value: strconv.FormatBool(parsed)})
		}
	}
	appendSingle := func(name operation.OptionName, value *operation.SingleValue) {
		if parsed, explicit := value.Value(); explicit {
			options = append(options, operation.Option{Name: name, Value: parsed})
		}
	}
	appendBool(operation.OptionArchive, &values.archive)
	appendBool(operation.OptionExtract, &values.extract)
	appendSingle(operation.OptionListen, &values.listen)
	appendBool(operation.OptionBackground, &values.background)
	appendSingle(operation.OptionAuth, &values.auth)
	appendSingle(operation.OptionAuthAttempts, &values.authAttempts)
	appendSingle(operation.OptionAuthFailAction, &values.authFailAction)
	appendSingle(operation.OptionLimit, &values.limit)
	appendBool(operation.OptionNoUI, &values.noUI)
	for _, value := range values.allowIP.Values() {
		options = append(options, operation.Option{Name: operation.OptionAllowIP, Value: value})
	}
	appendSingle(operation.OptionMaxFileSize, &values.maxFileSize)
	appendSingle(operation.OptionMaxExtractedSize, &values.maxExtractedSize)
	appendSingle(operation.OptionUploadRate, &values.uploadRate)
	appendSingle(operation.OptionDownloadRate, &values.downloadRate)
	options = append(options, values.selection.Options()...)
	return options
}

func newTransferCommand(dependencies Dependencies) *cobra.Command {
	values := &transferFlagValues{}
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
			plan, err := operation.Build(operation.Request{Source: args[0], Destination: args[2], Options: values.options()})
			if err != nil {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: err}
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
			return runRoute(command.Context(), dependencies, plan, selector, command.OutOrStdout(), command.ErrOrStderr())
		},
	}
	boolFlag(command, &values.archive, "archive", "create and transfer <source-name>.tar.gz")
	boolFlag(command, &values.extract, "extract", "extract a tar.gz archive into the destination root")
	command.Flags().Var(&values.listen, "listen", "incoming HTTP delivery bind address (default 127.0.0.1:8080)")
	boolFlag(command, &values.background, "background", "keep the delivery active after this command exits")
	command.Flags().Var(&values.auth, "auth", "authentication mode: none, basic, or password")
	command.Flags().Var(&values.authAttempts, "auth-attempts", "failed authentication threshold (default 5)")
	command.Flags().Var(&values.authFailAction, "auth-fail-action", "threshold action: ban or stop (default ban)")
	command.Flags().Var(&values.limit, "limit", "maximum concurrent transfers or unlimited")
	boolFlag(command, &values.noUI, "no-ui", "serve only the versioned data API")
	command.Flags().Var(&values.allowIP, "allow-ip", "allow a peer IP or CIDR (repeatable)")
	command.Flags().Var(values.selection.For(operation.OptionExclude), "exclude", "exclude a gitignore pattern")
	command.Flags().Var(values.selection.For(operation.OptionExcludeRegex), "exclude-regex", "exclude paths matching a Go regular expression")
	command.Flags().Var(values.selection.For(operation.OptionExcludeFrom), "exclude-from", "read gitignore patterns from a local file")
	command.Flags().Var(&values.maxFileSize, "max-file-size", "maximum incoming file size or unlimited (default 10GiB)")
	command.Flags().Var(&values.maxExtractedSize, "max-extracted-size", "maximum expanded size or unlimited (default 100GiB)")
	command.Flags().Var(&values.uploadRate, "upload-rate", "aggregate upload rate or unlimited")
	command.Flags().Var(&values.downloadRate, "download-rate", "aggregate download rate or unlimited")
	return command
}

func runRoute(ctx context.Context, dependencies Dependencies, plan operation.Plan, selector selection.Selector, stdout, stderr io.Writer) error {
	switch plan.Route {
	case operation.RoutePathToPath:
		return runTransfer(ctx, dependencies, plan, selector, stdout, stderr)
	case operation.RouteWebToPath, operation.RoutePathToWeb, operation.RouteWebhookToPath:
		if dependencies.Hosted == nil {
			return &commandError{code: ExitControl, stage: "control", cause: errors.New("web delivery dependencies are incomplete")}
		}
		if err := dependencies.Hosted(ctx, plan, stdout); err != nil {
			return &commandError{code: ExitControl, stage: "control", cause: err}
		}
		return nil
	case operation.RoutePathToHTTP:
		return runOutgoingWebhook(ctx, dependencies, plan, selector, stdout, stderr)
	default:
		return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("operation route has no runtime handler")}
	}
}

func boolFlag(command *cobra.Command, value *operation.BoolValue, name, usage string) {
	command.Flags().Var(value, name, usage)
	command.Flags().Lookup(name).NoOptDefVal = "true"
}
