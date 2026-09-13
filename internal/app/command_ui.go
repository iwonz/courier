package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/iwonz/courier/internal/admin"
	"github.com/iwonz/courier/internal/buildinfo"
	"github.com/iwonz/courier/internal/operation"
	"github.com/spf13/cobra"
)

const defaultAdminBind = "127.0.0.1:9090"

func adminManager() (admin.Manager, error) {
	directory, err := defaultStateDirectory()
	if err != nil {
		return admin.Manager{}, err
	}
	return admin.Manager{StateDirectory: directory, Compatibility: "admin-v1/" + buildinfo.Version}, nil
}

func startUI(ctx context.Context, request admin.StartRequest) (admin.StartResult, error) {
	manager, err := adminManager()
	if err != nil {
		return admin.StartResult{}, err
	}
	return manager.Start(ctx, request)
}

func stopUI(ctx context.Context) (admin.StopResult, error) {
	manager, err := adminManager()
	if err != nil {
		return admin.StopResult{}, err
	}
	return manager.Stop(ctx)
}

func newUICommand(start func(context.Context, admin.StartRequest) (admin.StartResult, error), stop func(context.Context) (admin.StopResult, error)) *cobra.Command {
	command := &cobra.Command{Use: "ui", Short: "Control the local Courier administration UI"}
	command.AddCommand(newUIStartCommand(start), newUIStopCommand(stop))
	return command
}

func newUIStartCommand(start func(context.Context, admin.StartRequest) (admin.StartResult, error)) *cobra.Command {
	var listen operation.SingleValue
	var background operation.BoolValue
	command := &cobra.Command{
		Use:   "start [options]",
		Short: "Start the local Courier administration UI",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			bind, set := listen.Value()
			if !set {
				bind = defaultAdminBind
			}
			if _, err := admin.CanonicalBind(bind); err != nil {
				return err
			}
			if start == nil {
				return controlCommandError(errors.New("administration UI startup is unavailable"))
			}
			backgroundValue, _ := background.Value()
			request := admin.StartRequest{Bind: bind, Background: backgroundValue}
			request.Ready = func(state admin.State) error {
				_, err := fmt.Fprintf(command.OutOrStdout(), "Courier administration UI: %s\n", admin.URL(state))
				return err
			}
			if _, err := start(command.Context(), request); err != nil {
				return controlCommandError(err)
			}
			return nil
		},
	}
	command.Flags().Var(&listen, "listen", "loopback administration address (default 127.0.0.1:9090)")
	boolFlag(command, &background, "background", "run the administration UI in the background")
	return command
}

func newUIStopCommand(stop func(context.Context) (admin.StopResult, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the local Courier administration UI",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if stop == nil {
				return controlCommandError(errors.New("administration UI stopping is unavailable"))
			}
			result, err := stop(command.Context())
			if err != nil {
				return controlCommandError(err)
			}
			message := "Stopped Courier administration UI.\n"
			if result.AlreadyStopped {
				message = "Courier administration UI was already stopped.\n"
			}
			if _, err := fmt.Fprint(command.OutOrStdout(), message); err != nil {
				return controlCommandError(err)
			}
			return nil
		},
	}
}
