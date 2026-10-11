package terminalui

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"rsc.io/qr"
)

func TestDetect(t *testing.T) {
	originalTerminal, originalSize := isTerminal, terminalSize
	t.Cleanup(func() { isTerminal, terminalSize = originalTerminal, originalSize })

	if mode := Detect(&bytes.Buffer{}, func(string) string { return "" }); mode.Interactive || mode.Color || mode.Width != defaultWidth {
		t.Fatalf("non-file mode=%+v", mode)
	}
	file, err := os.CreateTemp(t.TempDir(), "terminal")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = file.Close() })
	isTerminal = func(int) bool { return false }
	if mode := Detect(file, func(string) string { return "" }); mode.Interactive {
		t.Fatalf("non-terminal mode=%+v", mode)
	}

	isTerminal = func(int) bool { return true }
	terminalSize = func(int) (int, int, error) { return 144, 40, nil }
	if mode := Detect(file, func(string) string { return "" }); !mode.Interactive || !mode.Color || mode.Width != 144 {
		t.Fatalf("terminal mode=%+v", mode)
	}
	if mode := Detect(file, func(key string) string {
		if key == "NO_COLOR" {
			return "1"
		}
		return ""
	}); mode.Color {
		t.Fatalf("NO_COLOR mode=%+v", mode)
	}
	if mode := Detect(file, func(key string) string {
		if key == "TERM" {
			return "DuMb"
		}
		return ""
	}); mode.Interactive || mode.Color || mode.Width != 144 {
		t.Fatalf("dumb terminal mode=%+v", mode)
	}
	terminalSize = func(int) (int, int, error) { return 0, 0, os.ErrInvalid }
	if mode := Detect(file, func(string) string { return "" }); mode.Width != defaultWidth {
		t.Fatalf("fallback mode=%+v", mode)
	}
	terminalSize = func(int) (int, int, error) { return 0, 40, nil }
	if mode := Detect(file, func(string) string { return "" }); mode.Width != defaultWidth {
		t.Fatalf("zero-width mode=%+v", mode)
	}
}

func TestRendererPlainAndInteractive(t *testing.T) {
	var output bytes.Buffer
	plain := New(&output, Mode{Width: -1})
	if plain.Mode().Width != defaultWidth || plain.Text(ToneFailure, "bad\x1b[31m\x00") != "bad" {
		t.Fatalf("plain mode=%+v text=%q", plain.Mode(), plain.Text(ToneFailure, "bad\x1b[31m\x00"))
	}
	panel := plain.Panel("Title\x1b[2J", ToneInfo, []Field{{Label: "Path\x00", Value: "/tmp/a\x1b[31m"}}, "Footer\x07")
	if panel != "Title\nPath: /tmp/a\nFooter\n" {
		t.Fatalf("plain panel=%q", panel)
	}
	tableText := plain.Table("Rows", []string{"A\x00", "B"}, [][]string{{"one\x1b[2J", "two"}})
	if tableText != "Rows\nA\tB\none\ttwo\n" {
		t.Fatalf("plain table=%q", tableText)
	}

	interactive := New(&output, Mode{Interactive: true, Color: false, Width: 30})
	if got := interactive.Text(ToneWarning, "wait\x00"); got != "wait" {
		t.Fatalf("interactive text=%q", got)
	}
	for _, tone := range []Tone{ToneInfo, ToneSuccess, ToneWarning, ToneFailure} {
		got := interactive.Panel("State", tone, []Field{{Label: "A", Value: "B"}}, "done")
		if !strings.Contains(got, "STATE") || !strings.Contains(got, "A") || !strings.Contains(got, "done") {
			t.Fatalf("tone %d panel=%q", tone, got)
		}
	}
	withoutFooter := interactive.Panel("State", ToneInfo, nil, "")
	if !strings.Contains(withoutFooter, "STATE") {
		t.Fatalf("panel=%q", withoutFooter)
	}
	got := interactive.Table("Things", []string{"Name", "Value"}, [][]string{{"one", "two"}})
	if !strings.Contains(got, "THINGS") || !strings.Contains(got, "Name") || !strings.Contains(got, "one") {
		t.Fatalf("table=%q", got)
	}
}

func TestRendererUsesCompactThemeNeutralLayouts(t *testing.T) {
	var output bytes.Buffer
	wide := New(&output, Mode{Interactive: true, Color: true, Width: 72})
	panel := wide.Panel("Release notes", ToneInfo, []Field{
		{Value: "## Changelog\n* 8759741 feat: add browser upload queue progress"},
		{Label: "Release", Value: "https://github.com/iwonz/courier/releases/tag/v0.3.17"},
		{Label: "Reason", Value: "This deliberately long labeled value wraps onto a continuation line while preserving every word."},
	}, "Verified release metadata")
	assertThemeNeutral(t, panel)
	clean := ansi.Strip(panel)
	for _, wanted := range []string{"• RELEASE NOTES", "  ## Changelog", "  * 8759741", "  Release  https://github.com/", "  Reason   This deliberately long labeled value wraps onto a", "           continuation line while preserving every word.", "  Verified release metadata"} {
		if !strings.Contains(clean, wanted) {
			t.Fatalf("wide panel missing %q: %q", wanted, clean)
		}
	}
	complete := strings.NewReplacer(" ", "", "\n", "").Replace(clean)
	if !strings.Contains(complete, "Thisdeliberatelylonglabeledvaluewrapsontoacontinuationlinewhilepreservingeveryword.") {
		t.Fatalf("wide panel lost wrapped field content: %q", clean)
	}
	if longestLine(clean) >= 72 {
		t.Fatalf("wide panel paints terminal width: %q", clean)
	}

	narrow := New(&output, Mode{Interactive: true, Color: false, Width: 24}).Panel("Waiting", ToneWarning, []Field{{Label: "Destination", Value: "/a/very/long/destination/path"}}, "! Press Ctrl+C")
	clean = ansi.Strip(narrow)
	for _, wanted := range []string{"! WAITING", "  Destination", "    /a/very/long/", "  ! Press Ctrl+C"} {
		if !strings.Contains(clean, wanted) {
			t.Fatalf("narrow panel missing %q: %q", wanted, clean)
		}
	}
	if longestLine(clean) > 24 || strings.Contains(narrow, "\x1b[3") {
		t.Fatalf("narrow panel width/color: %q", narrow)
	}

	markers := map[string]bool{}
	colors := map[string]bool{}
	for _, tone := range []Tone{ToneInfo, ToneSuccess, ToneWarning, ToneFailure} {
		value := wide.Panel("State", tone, nil, "")
		assertThemeNeutral(t, value)
		markers[strings.Fields(ansi.Strip(value))[0]] = true
		colors[wide.Text(tone, "state")] = true
	}
	if len(markers) != 4 || len(colors) != 4 {
		t.Fatalf("semantic markers=%v colors=%v", markers, colors)
	}

	compact := wide.Table("Things", []string{"Name", "Value"}, [][]string{{"one", "two"}})
	assertThemeNeutral(t, compact)
	if clean = ansi.Strip(compact); !strings.Contains(clean, "• THINGS") || longestLine(clean) >= 72 {
		t.Fatalf("compact table=%q", clean)
	}
	narrowTable := New(&output, Mode{Interactive: true, Color: true, Width: 28}).Table("Long values", []string{"Name", "Value"}, [][]string{{"artifact", "https://example.test/a/very/long/path"}})
	assertThemeNeutral(t, narrowTable)
	clean = ansi.Strip(narrowTable)
	complete = strings.NewReplacer(" ", "", "\n", "", "\t", "", "│", "").Replace(clean)
	if !strings.Contains(complete, "https://example.test/a/very/long/path") || longestLine(clean) > 28 {
		t.Fatalf("narrow table=%q", clean)
	}
}

func assertThemeNeutral(t *testing.T, value string) {
	t.Helper()
	for _, forbidden := range []string{"\x1b[40m", "\x1b[41m", "\x1b[42m", "\x1b[43m", "\x1b[44m", "\x1b[45m", "\x1b[46m", "\x1b[47m", "\x1b[48;", "\x1b[38;2;"} {
		if strings.Contains(value, forbidden) {
			t.Fatalf("theme-neutral output contains %q: %q", forbidden, value)
		}
	}
}

func longestLine(value string) int {
	longest := 0
	for _, line := range strings.Split(value, "\n") {
		longest = max(longest, lipgloss.Width(line))
	}
	return longest
}

func TestSanitize(t *testing.T) {
	got := Sanitize("ok\nnext\tvalue\r\x00\x1b[31mred\x1b[0m")
	if got != "ok\nnext\tvaluered" {
		t.Fatalf("sanitize=%q", got)
	}
}

func TestQRRendering(t *testing.T) {
	original := encodeQR
	t.Cleanup(func() { encodeQR = original })
	payload := "https://127.0.0.1:8080/d/token/"
	var encoded string
	encodeQR = func(value string, level qr.Level) (*qr.Code, error) {
		encoded = value
		if level != qr.M {
			t.Fatalf("level=%v", level)
		}
		return &qr.Code{Bitmap: []byte{0xc0, 0xc0}, Size: 2, Stride: 1, Scale: 1}, nil
	}
	plain, err := New(ioDiscard{}, Mode{Interactive: true, Width: 10}).QR(payload)
	if err != nil || encoded != payload || strings.Contains(plain, "\x1b[") || !strings.ContainsAny(plain, "▀▄█") {
		t.Fatalf("plain=%q encoded=%q err=%v", plain, encoded, err)
	}
	color, err := New(ioDiscard{}, Mode{Interactive: true, Color: true, Width: 10}).QR(payload)
	if err != nil || !strings.Contains(color, "\x1b[40m") || !strings.Contains(color, "\x1b[47m") || !strings.HasSuffix(color, "\x1b[0m\n") {
		t.Fatalf("color=%q err=%v", color, err)
	}
	if narrow, err := New(ioDiscard{}, Mode{Interactive: true, Width: 9}).QR(payload); err != nil || narrow != "" {
		t.Fatalf("narrow=%q err=%v", narrow, err)
	}
	if redirected, err := New(ioDiscard{}, Mode{Width: 80}).QR(payload); err != nil || redirected != "" {
		t.Fatalf("redirected=%q err=%v", redirected, err)
	}
	if unsafe, err := New(ioDiscard{}, Mode{Interactive: true, Width: 80}).QR(payload + "\x1b[31m"); err == nil || unsafe != "" {
		t.Fatalf("unsafe=%q err=%v", unsafe, err)
	}
	encodeQR = func(string, qr.Level) (*qr.Code, error) { return nil, errors.New("encode") }
	if _, err := New(ioDiscard{}, Mode{Interactive: true, Width: 80}).QR(payload); err == nil {
		t.Fatal("encoding error ignored")
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(data []byte) (int, error) { return len(data), nil }

func TestQRCells(t *testing.T) {
	for _, test := range []struct {
		top, bottom bool
		plain       rune
		color       string
	}{
		{false, false, ' ', "\x1b[47m "},
		{true, false, '▀', "\x1b[30;47m▀"},
		{false, true, '▄', "\x1b[30;47m▄"},
		{true, true, '█', "\x1b[40m "},
	} {
		if got := plainQRCell(test.top, test.bottom); got != test.plain {
			t.Fatalf("plain %v/%v=%q", test.top, test.bottom, got)
		}
		if got := colorQRCell(test.top, test.bottom); got != test.color {
			t.Fatalf("color %v/%v=%q", test.top, test.bottom, got)
		}
	}
	code := &qr.Code{Bitmap: []byte{0x80}, Size: 1, Stride: 1, Scale: 1}
	if !qrBlack(code, 0, 0) || qrBlack(code, -1, 0) || qrBlack(code, 0, -1) || qrBlack(code, 1, 0) || qrBlack(code, 0, 1) {
		t.Fatal("QR bounds are incorrect")
	}
}
