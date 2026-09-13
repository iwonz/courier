// Package app assembles Courier commands and transfer dependencies.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"

	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/buildinfo"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/safety"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/transfer"
	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

const (
	ExitOK         = 0
	ExitCLI        = 2
	ExitConnection = 10
	ExitTransfer   = 20
	ExitUpdate     = 30
)

// Resource is an opened endpoint backend with its bounded path.
type Resource struct {
	Endpoint endpoint.Endpoint
	Backend  fsx.Backend
	Path     string
	Close    func() error
}

// Dependencies makes command orchestration locally testable.
type Dependencies struct {
	Open         func(context.Context, endpoint.Endpoint) (*Resource, error)
	OpenArtifact func(string) (fsx.Backend, string, func() error, error)
	Transfer     func(context.Context, transfer.Request) (transfer.Result, error)
	Archive      func(context.Context, fsx.Backend, string, string, string, progress.Sink) (*archive.Artifact, error)
	Update       func(context.Context) (update.Result, error)
	Reporter     func(io.Writer, bool) *report.Reporter
	Terminal     func(io.Writer) bool
	TempDir      string
	Build        BuildIdentity
}

// BuildIdentity is the version metadata rendered by the version command.
type BuildIdentity struct{ Version, Commit, Date string }

type commandError struct {
	code      int
	stage     string
	confirmed int64
	cause     error
}

func (e *commandError) Error() string { return e.cause.Error() }
func (e *commandError) Unwrap() error { return e.cause }

type resourceOpenError struct {
	remote bool
	role   string
	cause  error
}

func (e *resourceOpenError) Error() string { return fmt.Sprintf("open %s: %v", e.role, e.cause) }
func (e *resourceOpenError) Unwrap() error { return e.cause }

var (
	userHomeDirectory = os.UserHomeDir
	loadSSHConfig     = sshx.LoadConfig
	emptySSHConfig    = sshx.EmptyConfig
	currentUser       = user.Current
	environmentValue  = os.Getenv
	openRootedPath    = fsx.OpenRootedForPath
	openSSHConnection = func(ctx context.Context, factory sshx.Factory, host, username string) (*sshx.Connection, error) {
		return factory.Open(ctx, host, username)
	}
	detectSSHPlatform  = sshx.DetectPlatform
	terminalAttached   = term.IsTerminal
	readTerminalSecret = term.ReadPassword
)

// DefaultDependencies wires native libraries and secure user configuration.
func DefaultDependencies(input *os.File, promptOutput io.Writer) (Dependencies, error) {
	home, err := userHomeDirectory()
	if err != nil {
		return Dependencies{}, err
	}
	configuration, err := loadSSHConfig(filepath.Join(home, ".ssh", "config"), home)
	if errors.Is(err, os.ErrNotExist) {
		configuration = emptySSHConfig(home)
	} else if err != nil {
		return Dependencies{}, err
	}
	defaultUser := ""
	if current, userErr := currentUser(); userErr == nil {
		defaultUser = current.Username
	}
	factory := sshx.Factory{
		Config:      configuration,
		DefaultUser: defaultUser,
		KnownHosts:  filepath.Join(home, ".ssh", "known_hosts"),
		AgentSocket: environmentValue("SSH_AUTH_SOCK"),
		Prompt:      terminalPrompt(input, promptOutput),
	}
	open := func(ctx context.Context, value endpoint.Endpoint) (*Resource, error) {
		if !value.Remote {
			backend, relative, err := openRootedPath(value.Path)
			if err != nil {
				return nil, err
			}
			return &Resource{Endpoint: value, Backend: backend, Path: relative, Close: backend.Close}, nil
		}
		connection, err := openSSHConnection(ctx, factory, value.Host, value.User)
		if err != nil {
			return nil, err
		}
		if _, _, err := detectSSHPlatform(ctx, connection); err != nil {
			_ = connection.Close()
			return nil, err
		}
		effective := value
		effective.Host = connection.Target.Host
		effective.User = connection.Target.User
		return &Resource{Endpoint: effective, Backend: sshx.NewSFTPBackend(connection.SFTP), Path: value.Path, Close: connection.Close}, nil
	}
	updater := update.Updater{Repository: "iwonz/courier", Version: buildinfo.Version, Token: environmentValue("GITHUB_TOKEN")}
	return Dependencies{
		Open: open,
		OpenArtifact: func(name string) (fsx.Backend, string, func() error, error) {
			backend, relative, err := openRootedPath(name)
			if err != nil {
				return nil, "", nil, err
			}
			return backend, relative, backend.Close, nil
		},
		Transfer: (transfer.Engine{}).Run,
		Archive:  archive.Create,
		Update:   updater.Run,
		Reporter: report.New,
		Terminal: writerIsTerminal,
		Build:    BuildIdentity{Version: buildinfo.Version, Commit: buildinfo.Commit, Date: buildinfo.Date},
	}, nil
}

func terminalPrompt(input *os.File, output io.Writer) sshx.SecretPrompt {
	return func(label string) ([]byte, error) {
		if input == nil || !terminalAttached(int(input.Fd())) {
			return nil, errors.New("interactive terminal is required for password authentication")
		}
		fmt.Fprint(output, label)
		secret, err := readTerminalSecret(int(input.Fd()))
		fmt.Fprintln(output)
		return secret, err
	}
}

func writerIsTerminal(output io.Writer) bool {
	file, ok := output.(*os.File)
	return ok && terminalAttached(int(file.Fd()))
}

// NewRoot builds the Cobra command tree without global state.
func NewRoot(dependencies Dependencies) *cobra.Command {
	root := &cobra.Command{
		Use:           "courier",
		Short:         "Safely transfer files and directories",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	root.AddCommand(newTransferCommand(dependencies), newVersionCommand(dependencies.Build), newUpdateCommand(dependencies.Update))
	return root
}

func newTransferCommand(dependencies Dependencies) *cobra.Command {
	archiveMode := false
	command := &cobra.Command{
		Use:   "from <source> to <destination>",
		Short: "Transfer a file or directory",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 3 || args[1] != "to" {
				return &commandError{code: ExitCLI, stage: string(progress.StagePreflight), cause: errors.New("expected: courier from <source> to <destination> [flags]")}
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			return runTransfer(command.Context(), dependencies, args[0], args[2], archiveMode, command.OutOrStdout(), command.ErrOrStderr())
		},
	}
	command.Flags().BoolVar(&archiveMode, "archive", false, "create and transfer <source-name>.tar.gz")
	return command
}

func newVersionCommand(identity BuildIdentity) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build date",
		Args:  cobra.NoArgs,
		Run: func(command *cobra.Command, _ []string) {
			fmt.Fprintf(command.OutOrStdout(), "courier %s (commit %s, built %s)\n", identity.Version, identity.Commit, identity.Date)
		},
	}
}

func newUpdateCommand(run func(context.Context) (update.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Install the latest verified GitHub release",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if run == nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: errors.New("updater is unavailable")}
			}
			result, err := run(command.Context())
			if err != nil {
				return &commandError{code: ExitUpdate, stage: "update", cause: err}
			}
			if result.Current {
				fmt.Fprintf(command.OutOrStdout(), "courier %s is current\n", result.From)
			} else {
				fmt.Fprintf(command.OutOrStdout(), "updated courier from %s to %s\n", result.From, result.To)
				if result.Notes != "" {
					fmt.Fprintln(command.OutOrStdout(), result.Notes)
				}
			}
			return nil
		},
	}
}

// Execute runs a root with explicit arguments and stable exit-code mapping.
func Execute(ctx context.Context, root *cobra.Command, args []string, stdout, stderr io.Writer) int {
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	err := root.ExecuteContext(ctx)
	if err == nil {
		return ExitOK
	}
	var commandErr *commandError
	if errors.As(err, &commandErr) {
		report.Failure(stderr, commandErr.stage, commandErr.cause, commandErr.confirmed)
		return commandErr.code
	}
	report.Failure(stderr, string(progress.StagePreflight), err, 0)
	return ExitCLI
}

type transferOutcome struct {
	destination string
	result      transfer.Result
}

func runTransfer(ctx context.Context, dependencies Dependencies, sourceText, destinationText string, archiveMode bool, stdout, stderr io.Writer) error {
	outcome, err := performTransfer(ctx, dependencies, sourceText, destinationText, archiveMode, stderr)
	if err != nil {
		return err
	}
	report.Success(stdout, sourceText, outcome.destination, outcome.result.Bytes, outcome.result.Elapsed)
	return nil
}

func performTransfer(ctx context.Context, dependencies Dependencies, sourceText, destinationText string, archiveMode bool, stderr io.Writer) (outcome transferOutcome, resultErr error) {
	sourceEndpoint, err := endpoint.Parse(sourceText)
	if err != nil {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, err, 0)
	}
	destinationEndpoint, err := endpoint.Parse(destinationText)
	if err != nil {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, err, 0)
	}
	if dependencies.Open == nil || dependencies.Transfer == nil || dependencies.Reporter == nil || dependencies.Terminal == nil || archiveMode && (dependencies.Archive == nil || dependencies.OpenArtifact == nil) {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, errors.New("transfer dependencies are incomplete"), 0)
	}
	var source, destination *Resource
	group, groupContext := errgroup.WithContext(ctx)
	group.Go(func() error {
		opened, openErr := dependencies.Open(groupContext, sourceEndpoint)
		source = opened
		if openErr != nil {
			return &resourceOpenError{remote: sourceEndpoint.Remote, role: "source", cause: openErr}
		}
		return nil
	})
	group.Go(func() error {
		opened, openErr := dependencies.Open(groupContext, destinationEndpoint)
		destination = opened
		if openErr != nil {
			return &resourceOpenError{remote: destinationEndpoint.Remote, role: "destination", cause: openErr}
		}
		return nil
	})
	if err := group.Wait(); err != nil {
		_ = closeResources(source, destination)
		code := ExitTransfer
		stage := string(progress.StagePreflight)
		var openErr *resourceOpenError
		if errors.As(err, &openErr) && openErr.remote {
			code = ExitConnection
			stage = "connection"
		}
		return transferOutcome{}, &commandError{code: code, stage: stage, cause: err}
	}
	defer func() {
		if closeErr := closeResources(source, destination); closeErr != nil {
			resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, closeErr, 0))
		}
	}()
	sourceInfo, err := source.Backend.Lstat(source.Path)
	if err != nil {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, err, 0)
	}
	destinationIsDirectory := false
	if info, statErr := destination.Backend.Lstat(destination.Path); statErr == nil {
		destinationIsDirectory = info.IsDir()
	} else if !errors.Is(statErr, fs.ErrNotExist) {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, statErr, 0)
	}
	outputName := ""
	if archiveMode {
		outputName = source.Endpoint.Base() + ".tar.gz"
	}
	actualEndpoint, err := endpoint.ResolveDestination(source.Endpoint, destination.Endpoint, destinationIsDirectory, outputName)
	if err != nil {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, err, 0)
	}
	if err := safety.ValidateTransfer(source.Endpoint, actualEndpoint, sourceInfo.IsDir()); err != nil {
		return transferOutcome{}, transferCommandError(progress.StagePreflight, err, 0)
	}
	actualPath := destination.Path
	if outputName == "" {
		outputName = source.Endpoint.Base()
	}
	if actualEndpoint.Path != destination.Endpoint.Path {
		actualPath = destination.Backend.Join(destination.Path, outputName)
	}
	reporter := dependencies.Reporter(stderr, dependencies.Terminal(stderr))
	defer reporter.Finish()
	transferSource := source
	var artifact *archive.Artifact
	if archiveMode {
		artifact, err = dependencies.Archive(ctx, source.Backend, source.Path, source.Endpoint.Base(), dependencies.TempDir, reporter.Handle)
		if err != nil {
			return transferOutcome{}, transferCommandError(progress.StageArchive, err, 0)
		}
		defer func() {
			if cleanupErr := artifact.Cleanup(); cleanupErr != nil {
				resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, cleanupErr, 0))
			}
		}()
		archiveBackend, archivePath, closeArchive, openErr := dependencies.OpenArtifact(artifact.Path)
		if openErr != nil {
			return transferOutcome{}, transferCommandError(progress.StageArchive, openErr, 0)
		}
		defer func() {
			if closeArchive != nil {
				if closeErr := closeArchive(); closeErr != nil {
					resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, closeErr, 0))
				}
			}
		}()
		transferSource = &Resource{Backend: archiveBackend, Path: archivePath}
	}
	result, err := dependencies.Transfer(ctx, transfer.Request{SourceFS: transferSource.Backend, SourcePath: transferSource.Path, DestinationFS: destination.Backend, Destination: actualPath, Progress: reporter.Handle})
	if err != nil {
		var transferErr *transfer.Error
		if errors.As(err, &transferErr) {
			return transferOutcome{}, transferCommandError(transferErr.Stage, transferErr.Cause, transferErr.Confirmed)
		}
		return transferOutcome{}, transferCommandError(progress.StageTransfer, err, 0)
	}
	return transferOutcome{destination: actualEndpoint.Raw, result: result}, nil
}

func transferCommandError(stage progress.Stage, cause error, confirmed int64) *commandError {
	return &commandError{code: ExitTransfer, stage: string(stage), confirmed: confirmed, cause: cause}
}

func closeResources(resources ...*Resource) error {
	var result error
	for _, resource := range resources {
		if resource != nil && resource.Close != nil {
			result = errors.Join(result, resource.Close())
		}
	}
	return result
}
