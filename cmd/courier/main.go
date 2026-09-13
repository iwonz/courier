package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwonz/courier/internal/app"
)

var (
	dependencyFactory = app.DefaultDependencies
	exitProcess       = os.Exit
)

func run(ctx context.Context, args []string, input *os.File, stdout, stderr io.Writer) int {
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
