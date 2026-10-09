package terminalui

import (
	"bytes"
	"os"
	"strings"
	"testing"
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
	}); mode.Color {
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

func TestSanitize(t *testing.T) {
	got := Sanitize("ok\nnext\tvalue\r\x00\x1b[31mred\x1b[0m")
	if got != "ok\nnext\tvaluered" {
		t.Fatalf("sanitize=%q", got)
	}
}
