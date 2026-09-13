// Command contractdoc validates and generates Courier's CLI reference.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/iwonz/courier/internal/app"
	"github.com/iwonz/courier/internal/contract"
	"github.com/iwonz/courier/internal/operation"
	"github.com/spf13/cobra"
)

var (
	exitProcess = os.Exit
	readFile    = os.ReadFile
	writeFile   = os.WriteFile
	load        = contract.Load
	root        = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
)

const (
	contractPath  = "docs/cli-contract.yaml"
	referencePath = "docs/cli-reference.md"
)

func run(args []string, stdout, stderr io.Writer) int {
	mode := "--check"
	if len(args) > 1 || len(args) == 1 && args[0] != "--check" && args[0] != "--write" {
		fmt.Fprintln(stderr, "usage: contractdoc [--check|--write]")
		return 2
	}
	if len(args) == 1 {
		mode = args[0]
	}
	value, err := load(contractPath)
	if err != nil {
		fmt.Fprintf(stderr, "load CLI contract: %v\n", err)
		return 1
	}
	if err := value.CheckCobra(root()); err != nil {
		fmt.Fprintf(stderr, "check Cobra parity: %v\n", err)
		return 1
	}
	if err := value.CheckPlanner(operation.ContractMatrix()); err != nil {
		fmt.Fprintf(stderr, "check planner parity: %v\n", err)
		return 1
	}
	reference := value.Reference()
	if mode == "--write" {
		if err := writeFile(referencePath, reference, 0o644); err != nil {
			fmt.Fprintf(stderr, "write CLI reference: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, "generated", referencePath)
		return 0
	}
	current, err := readFile(referencePath)
	if err != nil {
		fmt.Fprintf(stderr, "read CLI reference: %v\n", err)
		return 1
	}
	if !bytes.Equal(current, reference) {
		fmt.Fprintln(stderr, "CLI reference is stale; run: go run ./cmd/contractdoc --write")
		return 1
	}
	fmt.Fprintln(stdout, "verified CLI contract and reference")
	return 0
}

func main() {
	exitProcess(run(os.Args[1:], os.Stdout, os.Stderr))
}
