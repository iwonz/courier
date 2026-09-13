package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwonz/courier/internal/buildinfo"
	"github.com/iwonz/courier/internal/cli"
)

const usage = `Courier safely transfers files and directories.

Usage:
  courier from <source> to <destination> [flags]
  courier <command> [flags]

Commands:`

var (
	newRegistry = cli.NewRegistry
	exitProcess = os.Exit
)

func versionCommand() cli.Command {
	return cli.CommandFunc{
		CommandName:     "version",
		CommandSynopsis: "print version, commit and build date",
		Execute: func(_ context.Context, args []string, stdout, _ io.Writer) error {
			if len(args) != 0 {
				return fmt.Errorf("%w: version takes no arguments", cli.ErrUsage)
			}
			fmt.Fprintf(stdout, "courier %s (commit %s, built %s)\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)
			return nil
		},
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	registry, err := newRegistry(usage, versionCommand())
	if err == nil {
		err = registry.Run(ctx, args, stdout, stderr)
	}
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, "courier:", err)
	if errors.Is(err, cli.ErrUsage) {
		return 2
	}
	return 1
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	exitProcess(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
