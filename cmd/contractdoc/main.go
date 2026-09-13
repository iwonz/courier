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
	landingPath   = "web/landing/src/contract.generated.json"
)

type generatedOutput struct {
	path string
	data []byte
}

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
	outputs := []generatedOutput{{path: referencePath, data: value.Reference()}, {path: landingPath, data: value.LandingJSON()}}
	if mode == "--write" {
		for _, output := range outputs {
			if err := writeFile(output.path, output.data, 0o644); err != nil {
				fmt.Fprintf(stderr, "write generated contract data %s: %v\n", output.path, err)
				return 1
			}
			fmt.Fprintln(stdout, "generated", output.path)
		}
		return 0
	}
	for _, output := range outputs {
		current, err := readFile(output.path)
		if err != nil {
			fmt.Fprintf(stderr, "read generated contract data %s: %v\n", output.path, err)
			return 1
		}
		if !bytes.Equal(current, output.data) {
			fmt.Fprintf(stderr, "%s is stale; run: go run ./cmd/contractdoc --write\n", output.path)
			return 1
		}
	}
	fmt.Fprintln(stdout, "verified CLI contract, reference, and landing data")
	return 0
}

func main() {
	exitProcess(run(os.Args[1:], os.Stdout, os.Stderr))
}
