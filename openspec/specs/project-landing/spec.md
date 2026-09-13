# project-landing Specification

## Purpose
Define the contract-backed, shared-kit Courier project site and its least-privilege publication as a repository-owned GitHub Pages artifact.

## Requirements

### Requirement: Shared-kit static landing

Courier SHALL provide a responsive static Lit application built from the shared UI kit with English fallback, Russian localization, accessible system/light/dark controls, reduced-motion support, and semantic keyboard navigation.

#### Scenario: A Russian-language browser visits for the first time

- **WHEN** the browser language resolves to Russian and no selector preference exists
- **THEN** the landing renders the Russian catalog while the same shipped contract data and installation commands remain available

### Requirement: Installation and release guidance

The landing SHALL present supported one-command installers, npm-compatible clients, Homebrew, Scoop, direct GitHub Releases, Linux package formats, shipped CLI routes, examples, documentation, source, license, and security links without embedding tokens or mutable release credentials.

#### Scenario: A user chooses a distribution channel

- **WHEN** the installation section is opened on any supported viewport
- **THEN** the exact repository-owned command or GitHub Release destination is visible and operable by keyboard

### Requirement: Repository-owned Pages deployment

Courier SHALL publish the deterministic landing artifact from `main` through GitHub Actions Pages deployment with least-privilege permissions, no external repository, and no `gh-pages` maintenance branch.

#### Scenario: Main passes the landing gate

- **WHEN** the Pages workflow builds a synchronized contract and tested production landing
- **THEN** it uploads and deploys only the static artifact through the GitHub Pages environment
