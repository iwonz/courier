# Context

The target surface includes planned server and browser workflows in addition to today's transfer, version, and update commands. The repository needs one source that distinguishes shipped behavior from planned behavior.

# Decisions

## YAML is the source of truth

`docs/cli-contract.yaml` records commands, endpoint kinds, routes, flags, applicability, conflicts, and unsupported forms. Generated Markdown and landing data are derived artifacts.

## Runtime remains native Go

The shipping binary does not parse documentation. A development tool validates YAML and compares shipped entries with the Cobra tree.

## Providers own commands

The root accepts command providers and validates unique command names before attachment. Cobra remains responsible for parsing and built-in help.

# Risks / Trade-offs

The contract duplicates a small amount of structural metadata from Cobra. A required parity check turns that duplication into an independently verified compatibility boundary.

# Migration Plan

Add the contract and generator, refactor root assembly, update the verification gate, regenerate the reference, and archive this delta into the baseline.
