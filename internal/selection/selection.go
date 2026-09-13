// Package selection compiles one ordered object-selection policy for every
// Courier traversal.
package selection

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/iwonz/courier/internal/operation"
)

const (
	maxRuleLineBytes = 64 * 1024
	maxRuleFileBytes = 1024 * 1024
)

var (
	// ErrInvalidRule identifies malformed or unsafe selection input.
	ErrInvalidRule = errors.New("invalid selection rule")
	// OpenFile is the dependency used for local --exclude-from files.
	OpenFile = defaultOpenFile
)

// Selector decides whether a root-relative object participates in traversal.
type Selector interface {
	Include(string, bool) bool
}

// Matcher is an immutable compiled selection policy.
type Matcher struct {
	gitignore gitignore.Matcher
	regex     []*regexp.Regexp
}

// Include reports whether a slash-separated path relative to the data root is
// selected. The root itself is always selected.
func (m Matcher) Include(relative string, directory bool) bool {
	if relative == "" || relative == "." {
		return true
	}
	clean := path.Clean(relative)
	if path.IsAbs(relative) || clean == ".." || strings.HasPrefix(clean, "../") {
		return false
	}
	for _, expression := range m.regex {
		if expression.MatchString(clean) {
			return false
		}
	}
	if m.gitignore == nil {
		return true
	}
	return !m.gitignore.Match(strings.Split(clean, "/"), directory)
}

// All returns a selector that includes every safe relative path.
func All() Matcher { return Matcher{} }

// Compile expands ordered rule occurrences and builds one immutable matcher.
func Compile(rules []operation.SelectionRule, open func(string) (io.ReadCloser, error)) (Matcher, error) {
	patterns := make([]gitignore.Pattern, 0)
	expressions := make([]*regexp.Regexp, 0)
	for _, rule := range rules {
		switch rule.Kind {
		case operation.SelectionGitignore:
			patterns = append(patterns, gitignore.ParsePattern(rule.Value, nil))
		case operation.SelectionRegex:
			expression, err := regexp.Compile(rule.Value)
			if err != nil {
				return Matcher{}, fmt.Errorf("%w at position %d: %v", ErrInvalidRule, rule.Position, err)
			}
			expressions = append(expressions, expression)
		case operation.SelectionFile:
			if open == nil {
				return Matcher{}, fmt.Errorf("%w at position %d: rule-file opener is unavailable", ErrInvalidRule, rule.Position)
			}
			values, err := readRuleFile(rule.Value, open)
			if err != nil {
				return Matcher{}, fmt.Errorf("%w at position %d: %v", ErrInvalidRule, rule.Position, err)
			}
			for _, value := range values {
				patterns = append(patterns, gitignore.ParsePattern(value, nil))
			}
		default:
			return Matcher{}, fmt.Errorf("%w at position %d: unknown kind %q", ErrInvalidRule, rule.Position, rule.Kind)
		}
	}
	result := Matcher{regex: expressions}
	if len(patterns) != 0 {
		result.gitignore = gitignore.NewMatcher(patterns)
	}
	return result, nil
}

func readRuleFile(name string, open func(string) (io.ReadCloser, error)) (values []string, resultErr error) {
	file, err := open(name)
	if err != nil {
		return nil, err
	}
	defer func() { resultErr = errors.Join(resultErr, file.Close()) }()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 4096), maxRuleLineBytes)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		total += len(line) + 1
		if total > maxRuleFileBytes {
			return nil, fmt.Errorf("rule file exceeds %d bytes", maxRuleFileBytes)
		}
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		values = append(values, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
