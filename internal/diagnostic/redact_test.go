package diagnostic

import "testing"

func TestRedact(t *testing.T) {
	tests := map[string]string{
		"plain failure": "plain failure",
		"GET https://user:pass@example.test/a?token=x#part": "GET https://example.test/a",
		"open ssh://user@example.test/path.":                "open ssh://example.test/path.",
		"password=correct horse reason denied":              "password=[REDACTED]",
		"Authorization: Bearer value\nnext":                 "Authorization=[REDACTED]\nnext",
		"cookie=session=x; other=y":                         "cookie=[REDACTED]",
		"bad:// url token":                                  "bad:// url token",
		"open https://example.test/%zz":                     "open https://example.test/%zz",
	}
	for input, want := range tests {
		if got := Redact(input); got != want {
			t.Errorf("Redact(%q)=%q want %q", input, got, want)
		}
	}
}
