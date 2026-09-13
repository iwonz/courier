// Package cli provides the top-level command registry.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ErrUsage identifies invalid command-line input.
var ErrUsage = errors.New("invalid command line")

// Command is an independently registered top-level CLI handler.
type Command interface {
	Name() string
	Synopsis() string
	Run(context.Context, []string, io.Writer, io.Writer) error
}

// Registry resolves and executes top-level commands.
type Registry struct {
	commands map[string]Command
	usage    string
}

// NewRegistry builds a registry and rejects duplicate or empty names.
func NewRegistry(usage string, commands ...Command) (*Registry, error) {
	r := &Registry{commands: make(map[string]Command), usage: usage}
	for _, command := range commands {
		if command == nil || strings.TrimSpace(command.Name()) == "" {
			return nil, fmt.Errorf("%w: command name is empty", ErrUsage)
		}
		if _, exists := r.commands[command.Name()]; exists {
			return nil, fmt.Errorf("%w: duplicate command %q", ErrUsage, command.Name())
		}
		r.commands[command.Name()] = command
	}
	return r, nil
}

// Run dispatches args or prints root help.
func (r *Registry) Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		r.Help(stdout)
		return nil
	}
	command, ok := r.commands[args[0]]
	if !ok {
		return fmt.Errorf("%w: unknown command %q", ErrUsage, args[0])
	}
	return command.Run(ctx, args[1:], stdout, stderr)
}

// Help writes deterministic root usage and command summaries.
func (r *Registry) Help(output io.Writer) {
	fmt.Fprintln(output, r.usage)
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(output, "  %-12s %s\n", name, r.commands[name].Synopsis())
	}
}

// CommandFunc adapts functions to Command.
type CommandFunc struct {
	CommandName     string
	CommandSynopsis string
	Execute         func(context.Context, []string, io.Writer, io.Writer) error
}

func (c CommandFunc) Name() string     { return c.CommandName }
func (c CommandFunc) Synopsis() string { return c.CommandSynopsis }
func (c CommandFunc) Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	return c.Execute(ctx, args, stdout, stderr)
}
