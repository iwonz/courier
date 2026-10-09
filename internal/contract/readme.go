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
		for _, argument := range command.Arguments {
			fmt.Fprintf(&output, "| `%s` | `<%s>` | %s | %s |\n", command.Path, argument.Name, requirement(requiredValue(argument.Required)), markdownCell(argument.Description.EN))
		}
	}
	output.WriteString("\n### Options\n\n| Option | Status | Description | Applies to | Default | Repeatable | Requires / conflicts |\n|---|---|---|---|---|---:|---|\n")
	for _, flag := range c.Flags {
		if flag.Status != "shipped" {
			continue
		}
		relations := "none"
		parts := make([]string, 0, 2)
		if len(flag.Requires) != 0 {
			parts = append(parts, "requires "+strings.Join(flag.Requires, ", "))
		}
		if len(flag.Conflicts) != 0 {
			parts = append(parts, "conflicts with "+strings.Join(flag.Conflicts, ", "))
		}
		if len(parts) != 0 {
			relations = strings.Join(parts, "; ")
		}
		fmt.Fprintf(&output, "| `%s` | %s | %s | %s | `%s` | %t | %s |\n", markdownCell(flag.Syntax), requirement(requiredValue(flag.Required)), markdownCell(flag.Description.EN), markdownCell(humanScopes(flag.AppliesTo, "en")), markdownCell(flag.Default), flag.Repeatable, relations)
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
