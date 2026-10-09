package contract

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	READMEStart = "<!-- BEGIN GENERATED CLI PARAMETERS -->"
	READMEEnd   = "<!-- END GENERATED CLI PARAMETERS -->"
)

// READMESection renders the managed public argument and option tables.
func (c Contract) READMESection() []byte {
	var output bytes.Buffer
	output.WriteString(READMEStart + "\n\n### Arguments\n\n| Command | Argument | Status | Description |\n|---|---|---|---|\n")
	for _, command := range c.Commands {
		if command.Status != "shipped" && command.Status != "system" {
			continue
		}
		for _, argument := range requiredArgumentsFirst(command.Arguments) {
			fmt.Fprintf(&output, "| `%s` | `<%s>` | %s | %s |\n", command.Path, argument.Name, requirement(requiredValue(argument.Required)), markdownCell(argument.Description.EN))
		}
	}
	output.WriteString("\n### Options\n\n| Option | Status | Description | Default | Repeatable | Applies to | Requires | Conflicts |\n|---|---|---|---|---:|---|---|---|\n")
	for _, flag := range requiredFlagsFirst(c.Flags) {
		if flag.Status != "shipped" {
			continue
		}
		requires := "none"
		if len(flag.Requires) != 0 {
			requires = strings.Join(flag.Requires, ", ")
		}
		conflicts := "none"
		if len(flag.Conflicts) != 0 {
			conflicts = strings.Join(flag.Conflicts, ", ")
		}
		fmt.Fprintf(&output, "| `%s` | %s | %s | `%s` | %t | %s | %s | %s |\n", markdownCell(flag.Syntax), requirement(requiredValue(flag.Required)), markdownCell(flag.Description.EN), markdownCell(flag.Default), flag.Repeatable, markdownCell(humanScopes(flag.AppliesTo, "en")), requires, conflicts)
	}
	output.WriteString("\n" + READMEEnd)
	return output.Bytes()
}

// UpdateREADME replaces exactly one stable managed section.
func (c Contract) UpdateREADME(current []byte) ([]byte, error) {
	start := bytes.Index(current, []byte(READMEStart))
	end := bytes.Index(current, []byte(READMEEnd))
	if start < 0 || end < 0 || end < start {
		return nil, errors.New("README CLI parameter markers are missing or out of order")
	}
	end += len(READMEEnd)
	if bytes.Index(current[start+len(READMEStart):], []byte(READMEStart)) >= 0 || bytes.Index(current[end:], []byte(READMEEnd)) >= 0 {
		return nil, errors.New("README CLI parameter markers are duplicated")
	}
	result := make([]byte, 0, len(current)+len(c.READMESection()))
	result = append(result, current[:start]...)
	result = append(result, c.READMESection()...)
	result = append(result, current[end:]...)
	return result, nil
}
