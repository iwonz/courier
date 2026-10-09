package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/iwonz/courier/internal/terminalui"
	"github.com/spf13/cobra"
)

func newUsageError(cause error) error { return &usageError{cause: cause} }

func renderUsageError(output io.Writer, mode terminalui.Mode, root *cobra.Command, args []string, cause error) {
	message, suggestion := usageMessage(cause)
	if suggestion == "" && strings.HasPrefix(strings.ToLower(message), "unknown command") && len(args) != 0 && !strings.HasPrefix(args[0], "-") {
		if candidates := root.SuggestionsFor(args[0]); len(candidates) != 0 {
			suggestion = candidates[0]
		}
	}
	renderer := terminalui.New(output, mode)
	_, _ = fmt.Fprintf(output, "%s: %s\n", renderer.Text(terminalui.ToneFailure, "error"), terminalui.Sanitize(message))
	if suggestion != "" {
		if !strings.HasPrefix(suggestion, root.Name()+" ") {
			suggestion = root.Name() + " " + suggestion
		}
		_, _ = fmt.Fprintf(output, "%s: %s\n", renderer.Text(terminalui.ToneWarning, "did you mean"), terminalui.Sanitize(suggestion))
	}
	_, _ = fmt.Fprintf(output, "%s: %s\n", renderer.Text(terminalui.ToneInfo, "help"), usageHelp(root, args))
}

func usageMessage(cause error) (string, string) {
	if cause == nil {
		return "invalid command usage", ""
	}
	lines := strings.Split(terminalui.Sanitize(cause.Error()), "\n")
	message := strings.TrimSpace(lines[0])
	suggestion := ""
	for index, line := range lines {
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "did you mean") {
			continue
		}
		for _, candidate := range lines[index+1:] {
			if candidate = strings.TrimSpace(candidate); candidate != "" {
				suggestion = candidate
				break
			}
			break
		}
	}
	if message == "" {
		message = "invalid command usage"
	}
	return message, suggestion
}

func usageHelp(root *cobra.Command, args []string) string {
	path := root.CommandPath()
	if command, _, err := root.Find(args); err == nil && command != nil {
		path = command.CommandPath()
	}
	return terminalui.Sanitize(path + " --help")
}
