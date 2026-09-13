package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/selection"
	"github.com/iwonz/courier/internal/webhook"
)

func runOutgoingWebhook(ctx context.Context, dependencies Dependencies, plan operation.Plan, selector selection.Selector, stdout, stderr io.Writer) (resultErr error) {
	if dependencies.Open == nil || dependencies.Webhook == nil || dependencies.DeliveryCredentials == nil || dependencies.Reporter == nil || dependencies.Terminal == nil || plan.Options.Archive && (dependencies.Archive == nil || dependencies.OpenArtifact == nil) {
		return transferCommandError(progress.StagePreflight, errors.New("outgoing webhook dependencies are incomplete"), 0)
	}
	credentials, err := dependencies.DeliveryCredentials(ctx, plan.Options.Auth)
	if err != nil {
		return transferCommandError(progress.StagePreflight, err, 0)
	}
	defer clearCredentials(&credentials)
	source, err := dependencies.Open(ctx, plan.Source)
	if err != nil {
		code := ExitTransfer
		stage := string(progress.StagePreflight)
		if plan.Source.Remote {
			code = ExitConnection
			stage = "connection"
		}
		return &commandError{code: code, stage: stage, cause: fmt.Errorf("open source: %w", err)}
	}
	defer func() {
		if source.Close != nil {
			if closeErr := source.Close(); closeErr != nil {
				resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, closeErr, 0))
			}
		}
	}()
	info, err := source.Backend.Lstat(source.Path)
	if err != nil {
		return transferCommandError(progress.StagePreflight, err, 0)
	}
	if info.IsDir() && !plan.Options.Archive {
		return transferCommandError(progress.StagePreflight, errors.New("outgoing webhook directory requires --archive"), 0)
	}
	if !plan.Options.Archive && !info.Mode().IsRegular() {
		return transferCommandError(progress.StagePreflight, errors.New("outgoing webhook source must be a regular file"), 0)
	}
	if !plan.Options.Archive && !selector.Include(plan.Source.Base(), false) {
		return transferCommandError(progress.StagePreflight, errors.New("outgoing webhook source is excluded"), 0)
	}

	reporter := dependencies.Reporter(stderr, dependencies.Terminal(stderr))
	defer reporter.Finish()
	sendBackend, sendPath, sendName := source.Backend, source.Path, plan.Source.Base()
	var artifact *archive.Artifact
	if plan.Options.Archive {
		artifact, err = dependencies.Archive(ctx, source.Backend, source.Path, plan.Source.Base(), dependencies.TempDir, selector, reporter.Handle)
		if err != nil {
			return transferCommandError(progress.StageArchive, err, 0)
		}
		defer func() {
			if cleanupErr := artifact.Cleanup(); cleanupErr != nil {
				resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, cleanupErr, 0))
			}
		}()
		archiveBackend, archivePath, closeArchive, openErr := dependencies.OpenArtifact(artifact.Path)
		if openErr != nil {
			return transferCommandError(progress.StageArchive, openErr, 0)
		}
		if closeArchive != nil {
			defer func() {
				if closeErr := closeArchive(); closeErr != nil {
					resultErr = errors.Join(resultErr, transferCommandError(progress.StageCleanup, closeErr, 0))
				}
			}()
		}
		sendBackend, sendPath, sendName = archiveBackend, archivePath, artifact.Name
	}
	result, err := dependencies.Webhook(ctx, webhook.Request{
		URL: plan.Destination.Raw, SourceFS: sendBackend, SourcePath: sendPath, Name: sendName,
		Username: credentials.BasicUsername, Password: credentials.BasicPassword,
		Rate: plan.Options.UploadRate.Value, Unlimited: plan.Options.UploadRate.Unlimited, Progress: reporter.Handle,
	})
	if err != nil {
		code := ExitTransfer
		var destinationError *url.Error
		if result.Bytes == 0 && errors.As(err, &destinationError) {
			code = ExitConnection
		}
		return &commandError{code: code, stage: string(progress.StageTransfer), read: result.Bytes, sent: result.Bytes, cause: err}
	}
	if _, err := fmt.Fprintf(stdout, "http-status: %d\n", result.StatusCode); err != nil {
		return transferCommandError(progress.StageComplete, err, result.Bytes)
	}
	report.Success(stdout, plan.Source.Raw, plan.Destination.Raw, result.Bytes, result.Elapsed)
	return nil
}
