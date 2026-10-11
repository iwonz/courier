// Package terminalui renders Courier's interactive terminal presentation.
package terminalui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"golang.org/x/term"
	"rsc.io/qr"
)

const defaultWidth = 100
const qrQuietZone = 4

var (
	isTerminal   = term.IsTerminal
	terminalSize = term.GetSize
	encodeQR     = qr.Encode
)

// Tone identifies a semantic terminal treatment.
type Tone uint8

const (
	ToneInfo Tone = iota
	ToneSuccess
	ToneWarning
	ToneFailure
)

// Field is one label/value row in a terminal panel.
type Field struct {
	Label string
	Value string
}

// Mode captures the output capabilities needed by the renderer.
type Mode struct {
	Interactive bool
	Color       bool
	Width       int
}

// Detect returns a terminal mode for one output writer and environment.
func Detect(output io.Writer, getenv func(string) string) Mode {
	mode := Mode{Width: defaultWidth}
	file, ok := output.(*os.File)
	if !ok || !isTerminal(int(file.Fd())) {
		return mode
	}
	mode.Interactive = true
	if width, _, err := terminalSize(int(file.Fd())); err == nil && width > 0 {
		mode.Width = width
	}
	if strings.EqualFold(getenv("TERM"), "dumb") {
		return Mode{Width: mode.Width}
	}
	mode.Color = getenv("NO_COLOR") == ""
	return mode
}

// Renderer owns the styles for one output stream.
type Renderer struct {
	mode     Mode
	renderer *lipgloss.Renderer
}

// New creates a renderer with an explicit, testable mode.
func New(output io.Writer, mode Mode) *Renderer {
	if mode.Width <= 0 {
		mode.Width = defaultWidth
	}
	renderer := lipgloss.NewRenderer(output)
	if mode.Color {
		renderer.SetColorProfile(termenv.ANSI)
	} else {
		renderer.SetColorProfile(termenv.Ascii)
	}
	return &Renderer{mode: mode, renderer: renderer}
}

// Mode returns the renderer capabilities.
func (r *Renderer) Mode() Mode { return r.mode }

// Text applies one semantic color without adding layout.
func (r *Renderer) Text(tone Tone, value string) string {
	value = Sanitize(value)
	if !r.mode.Interactive {
		return value
	}
	return r.renderer.NewStyle().Foreground(r.accent(tone)).Render(value)
}

// QR renders an exact text payload as a terminal QR code when the output can
// contain the complete code and its four-module quiet zone.
func (r *Renderer) QR(value string) (string, error) {
	if !r.mode.Interactive {
		return "", nil
	}
	clean := Sanitize(value)
	if clean == "" || clean != value {
		return "", errors.New("QR payload contains unsafe terminal text")
	}
	code, err := encodeQR(clean, qr.M)
	if err != nil {
		return "", err
	}
	size := code.Size + 2*qrQuietZone
	if size > r.mode.Width {
		return "", nil
	}
	var rendered strings.Builder
	for y := -qrQuietZone; y < code.Size+qrQuietZone; y += 2 {
		for x := -qrQuietZone; x < code.Size+qrQuietZone; x++ {
			top := qrBlack(code, x, y)
			bottom := qrBlack(code, x, y+1)
			if r.mode.Color {
				rendered.WriteString(colorQRCell(top, bottom))
			} else {
				rendered.WriteRune(plainQRCell(top, bottom))
			}
		}
		if r.mode.Color {
			rendered.WriteString("\x1b[0m")
		}
		rendered.WriteByte('\n')
	}
	return rendered.String(), nil
}

// Sanitize removes terminal control sequences while retaining readable layout.
func Sanitize(value string) string {
	value = ansi.Strip(value)
	return strings.Map(func(character rune) rune {
		if character == '\n' || character == '\t' || !unicode.IsControl(character) {
			return character
		}
		return -1
	}, value)
}

// Panel renders a titled semantic field group.
func (r *Renderer) Panel(title string, tone Tone, fields []Field, footer string) string {
	if !r.mode.Interactive {
		var plain strings.Builder
		fmt.Fprintln(&plain, Sanitize(title))
		for _, field := range fields {
			fmt.Fprintf(&plain, "%s: %s\n", Sanitize(field.Label), Sanitize(field.Value))
		}
		if footer != "" {
			fmt.Fprintln(&plain, Sanitize(footer))
		}
		return plain.String()
	}
	contentWidth := max(1, r.mode.Width-2)
	labelWidth := 0
	for _, field := range fields {
		if width := lipgloss.Width(Sanitize(field.Label)); width > labelWidth {
			labelWidth = width
		}
	}
	labelStyle := r.renderer.NewStyle().Bold(true).Width(labelWidth)
	lines := []string{r.heading(tone, title)}
	wide := r.mode.Width >= 48 && labelWidth+4 < contentWidth
	for _, field := range fields {
		label, value := Sanitize(field.Label), Sanitize(field.Value)
		if label == "" {
			for _, line := range wrapLines(value, contentWidth) {
				lines = append(lines, "  "+line)
			}
			continue
		}
		if !wide {
			lines = append(lines, "  "+r.renderer.NewStyle().Bold(true).Render(label))
			for _, line := range wrapLines(value, max(1, r.mode.Width-4)) {
				lines = append(lines, "    "+line)
			}
			continue
		}
		indent := strings.Repeat(" ", labelWidth+4)
		wrapped := wrapLines(value, max(1, contentWidth-labelWidth-2))
		lines = append(lines, "  "+labelStyle.Render(label)+"  "+wrapped[0])
		for _, line := range wrapped[1:] {
			lines = append(lines, indent+line)
		}
	}
	if footer != "" {
		lines = append(lines, "  "+r.renderer.NewStyle().Bold(true).Render(Sanitize(footer)))
	}
	return strings.Join(lines, "\n") + "\n\n"
}

// Table renders a responsive terminal table.
func (r *Renderer) Table(title string, headers []string, rows [][]string) string {
	cleanHeaders := cleanRow(headers)
	cleanRows := make([][]string, len(rows))
	for index, row := range rows {
		cleanRows[index] = cleanRow(row)
	}
	if !r.mode.Interactive {
		var plain strings.Builder
		fmt.Fprintln(&plain, Sanitize(title))
		fmt.Fprintln(&plain, strings.Join(cleanHeaders, "\t"))
		for _, row := range cleanRows {
			fmt.Fprintln(&plain, strings.Join(row, "\t"))
		}
		return plain.String()
	}
	base := r.renderer.NewStyle().Padding(0, 1)
	header := base.Bold(true)
	grid := table.New().
		Headers(cleanHeaders...).
		Rows(cleanRows...).
		Border(lipgloss.NormalBorder()).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return header
			}
			return base
		}).
		Wrap(true)
	rendered := grid.Render()
	if lipgloss.Width(rendered) > r.mode.Width {
		rendered = grid.Width(max(1, r.mode.Width)).Render()
	}
	return r.heading(ToneInfo, title) + "\n" + rendered + "\n\n"
}

func (r *Renderer) accent(tone Tone) lipgloss.Color {
	switch tone {
	case ToneSuccess:
		return lipgloss.Color("2")
	case ToneWarning:
		return lipgloss.Color("3")
	case ToneFailure:
		return lipgloss.Color("1")
	default:
		return lipgloss.Color("6")
	}
}

func (r *Renderer) heading(tone Tone, title string) string {
	marker := "•"
	switch tone {
	case ToneSuccess:
		marker = "✓"
	case ToneWarning:
		marker = "!"
	case ToneFailure:
		marker = "×"
	}
	marker = r.renderer.NewStyle().Foreground(r.accent(tone)).Render(marker)
	title = r.renderer.NewStyle().Bold(true).Render(strings.ToUpper(Sanitize(title)))
	return marker + " " + title
}

func wrapLines(value string, width int) []string {
	wrapped := ansi.Wrap(value, max(1, width), " /_")
	return strings.Split(wrapped, "\n")
}

func cleanRow(values []string) []string {
	clean := make([]string, len(values))
	for index, value := range values {
		clean[index] = Sanitize(value)
	}
	return clean
}

func qrBlack(code *qr.Code, x, y int) bool {
	return x >= 0 && y >= 0 && x < code.Size && y < code.Size && code.Black(x, y)
}

func plainQRCell(top, bottom bool) rune {
	switch {
	case top && bottom:
		return '█'
	case top:
		return '▀'
	case bottom:
		return '▄'
	default:
		return ' '
	}
}

func colorQRCell(top, bottom bool) string {
	switch {
	case top && bottom:
		return "\x1b[40m "
	case top:
		return "\x1b[30;47m▀"
	case bottom:
		return "\x1b[30;47m▄"
	default:
		return "\x1b[47m "
	}
}
