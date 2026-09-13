package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwonz/courier/internal/app"
	"github.com/iwonz/courier/internal/helper"
	"github.com/iwonz/courier/internal/update"
	"github.com/iwonz/courier/internal/worker"
)

var (
	dependencyFactory = app.DefaultDependencies
	exitProcess       = os.Exit
	runInternalUpdate = update.RunInternal
	serveHelperSFTP   = helper.ServeSFTP
	runInternalWorker = worker.RunInternal
)

func run(ctx context.Context, args []string, input *os.File, stdout, stderr io.Writer) int {
	if handled, err := runInternalUpdate(args); handled {
		if err != nil {
			fmt.Fprintf(stderr, "courier update handoff: %v\n", err)
			return app.ExitUpdate
		}
		return app.ExitOK
	}
	if len(args) == 1 && args[0] == "_helper-sftp" {
		if err := serveHelperSFTP(input, stdout); err != nil {
			fmt.Fprintf(stderr, "courier helper: %v\n", err)
			return app.ExitTransfer
		}
		return app.ExitOK
	}
	if handled, err := runInternalWorker(ctx, args); handled {
		if err != nil {
			fmt.Fprintf(stderr, "courier worker: %v\n", err)
			return app.ExitControl
		}
		return app.ExitOK
	}
	dependencies, err := dependencyFactory(input, stderr)
	if err != nil {
		fmt.Fprintln(stderr, "courier:", err)
		return app.ExitConnection
	}
	return app.Execute(ctx, app.NewRoot(dependencies), args, stdout, stderr)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	exitProcess(run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
