// Package diagnostic removes credential material from operator-visible text.
package diagnostic

import (
	"net/url"
	"regexp"
	"strings"
)

const replacement = "[REDACTED]"

var (
	urlPattern    = regexp.MustCompile(`[A-Za-z][A-Za-z0-9+.-]*://[^\s]+`)
	secretPattern = regexp.MustCompile(`(?i)\b(password|passwd|token|secret|authorization|proxy-authorization|cookie|set-cookie)\b\s*[:=]\s*[^\r\n]*`)
)

// Redact returns useful diagnostic context without URL credentials, query
// values, fragments, or conventional secret-bearing assignments.
func Redact(value string) string {
	value = urlPattern.ReplaceAllStringFunc(value, redactURL)
	return secretPattern.ReplaceAllString(value, `${1}=`+replacement)
}

func redactURL(value string) string {
	trailing := ""
	for len(value) != 0 && strings.ContainsRune(`.,;:!?)]}`, rune(value[len(value)-1])) {
		trailing = value[len(value)-1:] + trailing
		value = value[:len(value)-1]
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return value + trailing
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.ForceQuery = false
	parsed.Fragment = ""
	return parsed.String() + trailing
}
