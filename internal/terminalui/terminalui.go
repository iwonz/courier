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
	if !mode.Color {
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
	accent := r.accent(tone)
	labelWidth := 0
	for _, field := range fields {
		if width := lipgloss.Width(Sanitize(field.Label)); width > labelWidth {
			labelWidth = width
		}
	}
	labelStyle := r.renderer.NewStyle().Foreground(lipgloss.Color("#A8AFB8")).Width(labelWidth)
	valueStyle := r.renderer.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	lines := []string{r.renderer.NewStyle().Bold(true).Foreground(accent).Render(strings.ToUpper(Sanitize(title)))}
	for _, field := range fields {
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render(Sanitize(field.Label)), "  ", valueStyle.Render(Sanitize(field.Value))))
	}
	if footer != "" {
		lines = append(lines, r.renderer.NewStyle().Foreground(accent).Render(Sanitize(footer)))
	}
	width := r.mode.Width - 4
	if width < 36 {
		width = 36
	}
	return r.renderer.NewStyle().
		Width(width).
		Padding(1, 1).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0D1015")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(accent).
		Render(strings.Join(lines, "\n")) + "\n"
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
	base := r.renderer.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#0D1015")).Padding(0, 1)
	header := base.Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#71FFF6"))
	grid := table.New().
		Headers(cleanHeaders...).
		Rows(cleanRows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(r.renderer.NewStyle().Foreground(lipgloss.Color("#626A73"))).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return header
			}
			return base
		}).
		Wrap(true).
		Width(max(36, r.mode.Width-2))
	titleText := r.renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAD14F")).Render(strings.ToUpper(Sanitize(title)))
	return titleText + "\n" + grid.Render() + "\n"
}

func (r *Renderer) accent(tone Tone) lipgloss.Color {
	switch tone {
	case ToneSuccess:
		return lipgloss.Color("#71FFF6")
	case ToneWarning:
		return lipgloss.Color("#FAD14F")
	case ToneFailure:
		return lipgloss.Color("#C94A55")
	default:
		return lipgloss.Color("#71FFF6")
	}
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
