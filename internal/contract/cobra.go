package contract

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CheckCobra verifies every live command and local option against shipped contract entries.
func (c Contract) CheckCobra(root *cobra.Command) error {
	expectedCommands := map[string]Command{}
	expectedFlags := map[string]struct{}{}
	for _, command := range c.Commands {
		if command.Status == "shipped" || command.Status == "system" {
			expectedCommands[command.Path] = command
			for _, flag := range command.Flags {
				expectedFlags[command.Path+" --"+flag] = struct{}{}
			}
		}
	}
	actualCommands := map[string]*cobra.Command{}
	collectCommands(root, "", actualCommands)
	if err := equalKeys("command", expectedCommands, actualCommands); err != nil {
		return err
	}
	actualFlags := map[string]struct{}{}
	for path, command := range actualCommands {
		command.LocalNonPersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name != "help" {
				actualFlags[path+" --"+flag.Name] = struct{}{}
			}
		})
	}
	return equalStringKeys("flag", expectedFlags, actualFlags)
}

func collectCommands(parent *cobra.Command, prefix string, result map[string]*cobra.Command) {
	for _, command := range parent.Commands() {
		path := strings.TrimSpace(prefix + " " + command.Name())
		if command.Runnable() {
			result[path] = command
		}
		collectCommands(command, path, result)
	}
}

func equalKeys[T, U any](kind string, expected map[string]T, actual map[string]U) error {
	expectedSet := make(map[string]struct{}, len(expected))
	for key := range expected {
		expectedSet[key] = struct{}{}
	}
	actualSet := make(map[string]struct{}, len(actual))
	for key := range actual {
		actualSet[key] = struct{}{}
	}
	return equalStringKeys(kind, expectedSet, actualSet)
}

func equalStringKeys(kind string, expected, actual map[string]struct{}) error {
	missing := make([]string, 0)
	extra := make([]string, 0)
	for key := range expected {
		if _, ok := actual[key]; !ok {
			missing = append(missing, key)
		}
	}
	for key := range actual {
		if _, ok := expected[key]; !ok {
			extra = append(extra, key)
		}
	}
	if len(missing) == 0 && len(extra) == 0 {
		return nil
	}
	sort.Strings(missing)
	sort.Strings(extra)
	return fmt.Errorf("%s parity failed: missing=%v extra=%v", kind, missing, extra)
}
