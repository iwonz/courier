package app

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/operation"
	"github.com/spf13/cobra"
)

func listServers(ctx context.Context) ([]control.ServerView, error) {
	return withControlStore(ctx, func(service control.Service) ([]control.ServerView, error) {
		return service.List(ctx)
	})
}

func stopServers(ctx context.Context, request control.StopRequest) (control.StopResult, error) {
	return withControlStore(ctx, func(service control.Service) (control.StopResult, error) {
		return service.Stop(ctx, request)
	})
}

func withControlStore[T any](ctx context.Context, run func(control.Service) (T, error)) (result T, resultErr error) {
	directory, err := defaultStateDirectory()
	if err != nil {
		return result, err
	}
	store, err := openDeliveryStore(directory)
	if err != nil {
		return result, err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	return run(control.Service{Store: store})
}

func newServersCommand(list func(context.Context) ([]control.ServerView, error), stop func(context.Context, control.StopRequest) (control.StopResult, error)) *cobra.Command {
	command := &cobra.Command{
		Use:   "servers",
		Short: "List active and unreachable Courier data servers",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if list == nil {
				return controlCommandError(errors.New("server listing is unavailable"))
			}
			views, err := list(command.Context())
			if err != nil {
				return controlCommandError(err)
			}
			_, err = fmt.Fprint(command.OutOrStdout(), renderServers(views))
			if err != nil {
				return controlCommandError(err)
			}
			return nil
		},
	}
	command.AddCommand(newServersStopCommand(stop))
	return command
}

func newServersStopCommand(stop func(context.Context, control.StopRequest) (control.StopResult, error)) *cobra.Command {
	var all operation.BoolValue
	command := &cobra.Command{
		Use:   "stop <uuid>|--all",
		Short: "Stop one delivery, one data server, or all data servers",
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, args []string) error {
			allValue, _ := all.Value()
			request := control.StopRequest{All: allValue}
			switch {
			case allValue && len(args) != 0:
				return errors.New("--all cannot be combined with a UUID")
			case !allValue && len(args) != 1:
				return errors.New("expected one canonical UUID or --all")
			case !allValue:
				id, err := delivery.ParseID(args[0])
				if err != nil {
					return err
				}
				request.ID = id
			}
			if stop == nil {
				return controlCommandError(errors.New("server stopping is unavailable"))
			}
			result, err := stop(command.Context(), request)
			if outputErr := renderStopResult(command, request, result); outputErr != nil {
				return controlCommandError(errors.Join(err, outputErr))
			}
			if err != nil {
				return controlCommandError(err)
			}
			return nil
		},
	}
	boolFlag(command, &all, "all", "stop all Courier data servers")
	return command
}

func controlCommandError(err error) error {
	return &commandError{code: ExitControl, stage: "control", cause: err}
}

func renderServers(views []control.ServerView) string {
	if len(views) == 0 {
		return "No Courier data servers found.\n"
	}
	var output strings.Builder
	for _, view := range views {
		status := "unreachable"
		if view.Live {
			status = "live"
		}
		fmt.Fprintf(&output, "server %s\n", view.Server.ID)
		fmt.Fprintf(&output, "  status: %s\n", status)
		fmt.Fprintf(&output, "  bind: %s\n", view.Server.Bind)
		fmt.Fprintf(&output, "  process: %d\n", view.Server.ProcessID)
		fmt.Fprintf(&output, "  state: %s\n", view.Server.State)
		fmt.Fprintf(&output, "  started: %s\n", formatControlTime(view.Server.StartedAt))
		fmt.Fprintf(&output, "  updated: %s\n", formatControlTime(view.Server.UpdatedAt))
		for _, item := range view.Deliveries {
			fmt.Fprintf(&output, "  delivery %s\n", item.ID)
			fmt.Fprintf(&output, "    route: %s\n", item.Route)
			fmt.Fprintf(&output, "    source: %s\n", displayEndpoint(item.Source))
			fmt.Fprintf(&output, "    destination: %s\n", displayEndpoint(item.Destination))
			fmt.Fprintf(&output, "    state: %s\n", item.State)
			fmt.Fprintf(&output, "    policy: auth=%s attempts=%d fail-action=%s limit=%s max-file-size=%s max-extracted-size=%s upload-rate=%s download-rate=%s no-ui=%t allow-ip=%s\n",
				item.Policy.Auth, item.Policy.AuthAttempts, item.Policy.AuthFailAction, formatLimit(item.Policy.DeliveryLimit),
				formatLimit(item.Policy.MaxFileSize), formatLimit(item.Policy.MaxExtractedSize), formatLimit(item.Policy.UploadRate),
				formatLimit(item.Policy.DownloadRate), item.Policy.NoUI, formatAllowIP(item.Policy.AllowIP))
			fmt.Fprintf(&output, "    counters: read=%d sent=%d confirmed=%d\n", item.Counters.Read, item.Counters.Sent, item.Counters.Confirmed)
			fmt.Fprintf(&output, "    created: %s\n", formatControlTime(item.CreatedAt))
			fmt.Fprintf(&output, "    updated: %s\n", formatControlTime(item.UpdatedAt))
		}
	}
	return output.String()
}

func renderStopResult(command *cobra.Command, request control.StopRequest, result control.StopResult) error {
	switch {
	case request.All:
		_, err := fmt.Fprintf(command.OutOrStdout(), "Stopped data servers: %d\n", result.StoppedServers)
		return err
	case result.AlreadyStopped:
		_, err := fmt.Fprintf(command.OutOrStdout(), "%s %s was already stopped.\n", result.Kind, result.ID)
		return err
	case result.Kind != "":
		_, err := fmt.Fprintf(command.OutOrStdout(), "Stopped %s %s.\n", result.Kind, result.ID)
		return err
	default:
		return nil
	}
}

func displayEndpoint(value string) string {
	if value == "" {
		return "<unavailable>"
	}
	return value
}

func formatLimit(limit delivery.Limit) string {
	if limit.Unlimited {
		return "unlimited"
	}
	return fmt.Sprintf("%d", limit.Value)
}

func formatAllowIP(values []string) string {
	if len(values) == 0 {
		return "none"
	}
	copyOfValues := append([]string(nil), values...)
	sort.Strings(copyOfValues)
	return strings.Join(copyOfValues, ",")
}

func formatControlTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}
