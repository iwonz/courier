package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/iwonz/courier/internal/terminalui"
	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCommand(run func(context.Context, update.Sink) (update.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Install the latest verified GitHub release",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if run == nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: errors.New("updater is unavailable")}
			}
			progress := newUpdateProgress(command.ErrOrStderr(), terminalMode(command.ErrOrStderr()))
			result, err := run(command.Context(), progress.Handle)
			progress.Finish()
			if err != nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: err}
			}
			if result.ReleaseURL == "" {
				result.ReleaseURL = update.ReleaseURL("", result.To)
			}
			renderUpdateResult(command.OutOrStdout(), terminalMode(command.OutOrStdout()), result)
			return nil
		},
	}
}

type updateProgress struct {
	output io.Writer
	mode   terminalui.Mode
	live   bool
}

func newUpdateProgress(output io.Writer, mode terminalui.Mode) *updateProgress {
	return &updateProgress{output: output, mode: mode}
}

func (p *updateProgress) Handle(event update.Event) {
	if !p.mode.Interactive || event.Stage == update.StageComplete {
		return
	}
	label := updateStageLabel(event.Stage)
	marker, tone := "●", terminalui.ToneWarning
	if event.Complete {
		marker, tone = "✓", terminalui.ToneSuccess
	}
	detail := ""
	if event.Stage == update.StageDownloadArchive || event.Stage == update.StageDownloadChecksum {
		if event.Total > 0 {
			percentage := event.Current * 100 / event.Total
			if percentage > 100 {
				percentage = 100
			}
			detail = fmt.Sprintf(" · %d/%d bytes · %d%%", event.Current, event.Total, percentage)
		} else {
			detail = fmt.Sprintf(" · %d bytes", event.Current)
		}
	}
	line := terminalui.New(p.output, p.mode).Text(tone, marker+" "+label+detail)
	ending := ""
	if event.Complete {
		ending = "\n"
	}
	fmt.Fprintf(p.output, "\r\x1b[2K%s%s", line, ending)
	p.live = !event.Complete
}

func (p *updateProgress) Finish() {
	if p.mode.Interactive && p.live {
		fmt.Fprint(p.output, "\r\x1b[2K")
	}
	p.live = false
}

func updateStageLabel(stage update.Stage) string {
	switch stage {
	case update.StageCheck:
		return "Checking the latest release"
	case update.StageDownloadArchive:
		return "Downloading the release archive"
	case update.StageDownloadChecksum:
		return "Downloading checksums"
	case update.StageVerify:
		return "Verifying the checksum"
	case update.StageExtract:
		return "Extracting the Courier binary"
	case update.StageInstall:
		return "Installing the verified binary"
	default:
		return "Completing the update"
	}
}

func renderUpdateResult(output io.Writer, mode terminalui.Mode, result update.Result) {
	if !mode.Interactive {
		if result.Current {
			fmt.Fprintf(output, "courier %s is current\n", result.From)
		} else {
			fmt.Fprintf(output, "updated courier from %s to %s\n", result.From, result.To)
			if result.Notes != "" {
				fmt.Fprintln(output, result.Notes)
			}
		}
		fmt.Fprintf(output, "release: %s\n", result.ReleaseURL)
		return
	}
	renderer := terminalui.New(output, mode)
	if strings.TrimSpace(result.Notes) != "" && !result.Current {
		fmt.Fprint(output, renderer.Panel("Release notes", terminalui.ToneInfo, []terminalui.Field{{Label: "Notes", Value: result.Notes}}, "Verified release metadata"))
	}
	title := "Courier is up to date"
	fields := []terminalui.Field{{Label: "Version", Value: result.From}, {Label: "Release", Value: result.ReleaseURL}}
	if !result.Current {
		title = "Courier updated"
		fields = []terminalui.Field{{Label: "Previous", Value: result.From}, {Label: "Version", Value: result.To}, {Label: "Release", Value: result.ReleaseURL}}
	}
	fmt.Fprint(output, renderer.Panel(title, terminalui.ToneSuccess, fields, "✓ Verified release installed"))
}
