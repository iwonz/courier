package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwonz/courier/internal/admin"
	"github.com/iwonz/courier/internal/app"
	"github.com/iwonz/courier/internal/helper"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/update"
	"github.com/iwonz/courier/internal/worker"
)

var (
	dependencyFactory = app.DefaultDependencies
	exitProcess       = os.Exit
	runInternalUpdate = update.RunInternal
	runInternalAdmin  = admin.RunInternal
	serveHelperSFTP   = helper.ServeSFTP
	runInternalWorker = worker.RunInternal
)

func run(ctx context.Context, args []string, input *os.File, stdout, stderr io.Writer) int {
	if handled, err := runInternalUpdate(args); handled {
		if err != nil {
			return internalFailure(stderr, app.ExitUpdate, "update", err)
		}
		return app.ExitOK
	}
	if handled, err := runInternalAdmin(ctx, args); handled {
		if err != nil {
			return internalFailure(stderr, app.ExitControl, "control", err)
		}
		return app.ExitOK
	}
	if len(args) == 1 && args[0] == "_helper-sftp" {
		if err := serveHelperSFTP(input, stdout); err != nil {
			return internalFailure(stderr, app.ExitTransfer, "transfer", err)
		}
		return app.ExitOK
	}
	if handled, err := runInternalWorker(ctx, args); handled {
		if err != nil {
			return internalFailure(stderr, app.ExitControl, "control", err)
		}
		return app.ExitOK
	}
	dependencies, err := dependencyFactory(input, stderr)
	if err != nil {
		return internalFailure(stderr, app.ExitConnection, "connection", err)
	}
	return app.Execute(ctx, app.NewRoot(dependencies), args, stdout, stderr)
}

func internalFailure(stderr io.Writer, code int, stage string, err error) int {
	report.FailureCounters(stderr, stage, err, 0, 0, 0)
	if errors.Is(err, context.Canceled) {
		return app.ExitInterrupted
	}
	return code
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	exitProcess(run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
