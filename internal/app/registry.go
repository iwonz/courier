package app

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// CommandProvider constructs one isolated top-level command.
type CommandProvider interface {
	Command() *cobra.Command
}

// ProviderFunc adapts a command constructor to CommandProvider.
type ProviderFunc func() *cobra.Command

func (provider ProviderFunc) Command() *cobra.Command { return provider() }

// NewRootWithProviders validates providers and constructs the public Cobra tree.
func NewRootWithProviders(providers ...CommandProvider) (*cobra.Command, error) {
	root := &cobra.Command{
		Use:           "courier",
		Short:         "Safely transfer files and directories",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	seen := map[string]struct{}{}
	for _, provider := range providers {
		if provider == nil {
			return nil, errors.New("nil command provider")
		}
		command := provider.Command()
		if command == nil || command.Name() == "" {
			return nil, errors.New("command provider returned an invalid command")
		}
		if _, exists := seen[command.Name()]; exists {
			return nil, fmt.Errorf("duplicate command %q", command.Name())
		}
		seen[command.Name()] = struct{}{}
		root.AddCommand(command)
	}
	root.InitDefaultHelpCmd()
	return root, nil
}
