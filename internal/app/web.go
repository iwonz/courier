package app

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/buildinfo"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/ipc"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/terminalui"
	"github.com/iwonz/courier/internal/webdelivery"
	"github.com/iwonz/courier/internal/worker"
	"github.com/pkg/sftp"
)

var (
	defaultStateDirectory = delivery.DefaultStateDirectory
	openDeliveryStore     = delivery.OpenStore
	newWebDefinition      = webdelivery.NewDefinition
	webAddress            = webdelivery.URL
	webhookAddress        = webdelivery.WebhookURL
	webNow                = time.Now
	webStatusTicker       = func() (<-chan time.Time, func()) {
		ticker := time.NewTicker(time.Second)
		return ticker.C, ticker.Stop
	}
	newWebCoordinator = func(store *delivery.Store, directory string) webCoordinator {
		return worker.DefaultCoordinator(store, directory)
	}
	runWebAcquire   = acquireWebDelivery
	releaseWebLease = (*worker.Lease).Release
)

type webCoordinator interface {
	Acquire(context.Context, worker.AcquireRequest) (worker.Acquired, error)
}

func acquireWebDelivery(ctx context.Context, plan operation.Plan, configured delivery.Policy, runtimeDefinition json.RawMessage) (worker.Acquired, error) {
	stateDirectory, err := defaultStateDirectory()
	if err != nil {
		return worker.Acquired{}, err
	}
	store, err := openDeliveryStore(stateDirectory)
	if err != nil {
		return worker.Acquired{}, err
	}
	defer store.Close()
	route := delivery.RouteWebToPath
	if plan.Route == operation.RoutePathToWeb {
		route = delivery.RoutePathToWeb
	} else if plan.Route == operation.RouteWebhookToPath {
		route = delivery.RouteWebhookToPath
	}
	return newWebCoordinator(store, stateDirectory).Acquire(ctx, worker.AcquireRequest{
		Bind: plan.Options.Listen, Compatibility: "web-v1/" + buildinfo.Version, Route: route,
		Source: plan.Source.Raw, Destination: plan.Destination.Raw,
		Policy: configured, Foreground: !plan.Options.Background, At: webNow().UTC(), RuntimeDefinition: runtimeDefinition,
	})
}

func deliveryCredentialPrompt(input *os.File, output io.Writer) func(context.Context, operation.AuthMode) (policy.Credentials, error) {
	return func(ctx context.Context, mode operation.AuthMode) (policy.Credentials, error) {
		if mode == operation.AuthNone {
			return policy.Credentials{}, nil
		}
		if input == nil || !terminalAttached(int(input.Fd())) {
			return policy.Credentials{}, errors.New("interactive terminal is required for delivery credentials")
		}
		credentials := policy.Credentials{}
		if mode == operation.AuthBasic {
			label := terminalui.New(output, terminalMode(output)).Text(terminalui.ToneWarning, "Basic username: ")
			if _, err := fmt.Fprint(output, label); err != nil {
				return policy.Credentials{}, err
			}
			username, err := bufio.NewReader(input).ReadString('\n')
			if err != nil {
				return policy.Credentials{}, err
			}
			credentials.BasicUsername = strings.TrimSpace(username)
			if credentials.BasicUsername == "" {
				return policy.Credentials{}, errors.New("Basic username must not be empty")
			}
		}
		if err := ctx.Err(); err != nil {
			return policy.Credentials{}, err
		}
		secret, err := terminalPrompt(input, output)("Delivery password")
		if err != nil {
			return policy.Credentials{}, err
		}
		if mode == operation.AuthBasic {
			credentials.BasicPassword = secret
		} else {
			credentials.Password = secret
		}
		return credentials, nil
	}
}

// HostedEndpoint carries an opened endpoint for preflight and the bounded
// credentials needed for the worker to reopen it after acquisition.
type HostedEndpoint struct {
	webdelivery.EndpointRuntime
	Resource *Resource
}

func endpointCredentialProvider(factory sshx.Factory) func(context.Context, endpoint.Endpoint) (*HostedEndpoint, error) {
	return func(ctx context.Context, value endpoint.Endpoint) (prepared *HostedEndpoint, resultErr error) {
		if !value.Remote {
			backend, relative, err := openRootedPath(value.Path)
			if err != nil {
				return nil, err
			}
			return &HostedEndpoint{Resource: &Resource{Endpoint: value, Backend: backend, Path: relative, Close: backend.Close}}, nil
		}
		runtime := webdelivery.EndpointRuntime{}
		var connection *sshx.Connection
		defer func() {
			if resultErr != nil {
				runtime.Clear()
				if connection != nil {
					_ = closeSSHConnection(connection)
				}
			}
		}()
		connectionFactory := factory
		prompt := factory.Prompt
		connectionFactory.Prompt = func(label string) ([]byte, error) {
			if prompt == nil {
				return nil, errors.New("interactive SSH credential prompt is unavailable")
			}
			secret, err := prompt(label)
			if err != nil {
				return nil, err
			}
			runtime.Credentials = append(runtime.Credentials, webdelivery.EndpointCredential{Prompt: label, Secret: append([]byte(nil), secret...)})
			return secret, nil
		}
		fallback := factory.SFTPFallback
		if fallback != nil {
			connectionFactory.SFTPFallback = func(ctx context.Context, connection *sshx.Connection, capability *sshx.CapabilityError) (*sftp.Client, io.Closer, error) {
				client, closer, err := fallback(ctx, connection, capability)
				if err == nil {
					runtime.AllowHelper = true
				}
				return client, closer, err
			}
		}
		connection, err := openSSHConnection(ctx, connectionFactory, value.Host, value.User)
		if err != nil {
			return nil, err
		}
		if _, _, err := detectSSHPlatform(ctx, connection); err != nil {
			return nil, err
		}
		effective := value
		effective.Host = connection.Target.Host
		effective.User = connection.Target.User
		resource := &Resource{Endpoint: effective, Backend: sshx.NewSFTPBackend(connection.SFTP), Path: value.Path, Close: func() error { return closeSSHConnection(connection) }}
		return &HostedEndpoint{EndpointRuntime: runtime, Resource: resource}, nil
	}
}

func webRunner(credentials func(context.Context, operation.AuthMode) (policy.Credentials, error), endpointCredentials func(context.Context, endpoint.Endpoint) (*HostedEndpoint, error), confirm DirectoryConfirm) func(context.Context, operation.Plan, io.Writer) error {
	return func(ctx context.Context, plan operation.Plan, output io.Writer) error {
		if credentials == nil {
			return errors.New("web credential provider is unavailable")
		}
		pathEndpoint := plan.Source
		role := "Source"
		expectation := pathFlexibleSource
		if plan.Route == operation.RouteWebToPath || plan.Route == operation.RouteWebhookToPath {
			pathEndpoint = plan.Destination
			role = "Destination"
			expectation = pathDirectory
		}
		if endpointCredentials == nil {
			return errors.New("web endpoint provider is unavailable")
		}
		prepared, err := endpointCredentials(ctx, pathEndpoint)
		if err != nil {
			return err
		}
		if prepared == nil || prepared.Resource == nil {
			return errors.New("web endpoint provider returned an incomplete resource")
		}
		endpointRuntime := prepared.EndpointRuntime
		defer endpointRuntime.Clear()
		_, preflightErr := preflightSinglePath(ctx, confirm, plan.Options.ForceSourceCreation, prepared.Resource, role, expectation)
		closeErr := closeResources(prepared.Resource)
		if preflightErr != nil || closeErr != nil {
			return errors.Join(preflightErr, closeErr)
		}
		secret, err := credentials(ctx, plan.Options.Auth)
		if err != nil {
			return err
		}
		defer clearCredentials(&secret)
		definition, err := newWebDefinition(plan, secret, endpointRuntime, rand.Reader)
		if err != nil {
			return err
		}
		defer definition.ClearSecrets()
		runtimeDefinition, _ := json.Marshal(definition)
		defer clear(runtimeDefinition)
		configured, err := deliveryPolicy(plan)
		if err != nil {
			return err
		}
		acquired, err := runWebAcquire(ctx, plan, configured, runtimeDefinition)
		if err != nil {
			return err
		}
		addressFunc := webAddress
		if plan.Route == operation.RouteWebhookToPath {
			addressFunc = webhookAddress
		}
		address, err := addressFunc(plan.Options.Listen, definition.Token)
		if err != nil {
			if acquired.Lease != nil {
				_ = releaseWebLease(acquired.Lease, context.Background())
			}
			return err
		}
		if err := renderHostedReady(output, plan, address, acquired.DeliveryID); err != nil {
			if acquired.Lease != nil {
				_ = releaseWebLease(acquired.Lease, context.Background())
			}
			return err
		}
		if plan.Options.Background {
			return nil
		}
		if acquired.Lease == nil {
			return errors.New("foreground delivery did not return a lease")
		}
		waitErr := waitForHostedStop(ctx, output, plan.Route, webNow())
		releaseContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := releaseWebLease(acquired.Lease, releaseContext); err != nil {
			return errors.Join(waitErr, err)
		}
		return waitErr
	}
}

func renderHostedReady(output io.Writer, plan operation.Plan, address string, id delivery.ID) error {
	mode := terminalMode(output)
	if !mode.Interactive {
		_, err := fmt.Fprintf(output, "delivery: %s\nid: %s\n", address, id)
		return err
	}
	runMode, footer := "foreground", "Press Ctrl+C to stop"
	if plan.Options.Background {
		runMode = "background"
		footer = "Stop with: courier servers stop " + string(id)
	}
	panel := terminalui.New(output, mode).Panel("Delivery ready", terminalui.ToneSuccess, []terminalui.Field{
		{Label: "Route", Value: plan.Route.String()},
		{Label: "URL", Value: address},
		{Label: "ID", Value: string(id)},
		{Label: "Source", Value: plan.Source.Raw},
		{Label: "Destination", Value: plan.Destination.Raw},
		{Label: "Mode", Value: runMode},
	}, footer)
	_, err := fmt.Fprint(output, panel)
	return err
}

func waitForHostedStop(ctx context.Context, output io.Writer, route operation.Route, started time.Time) error {
	mode := terminalMode(output)
	if !mode.Interactive {
		<-ctx.Done()
		return ctx.Err()
	}
	status := "Serving browser downloads"
	if route == operation.RouteWebToPath {
		status = "Waiting for a browser upload"
	} else if route == operation.RouteWebhookToPath {
		status = "Listening for webhook uploads"
	}
	ticks, stop := webStatusTicker()
	defer stop()
	renderer := terminalui.New(output, mode)
	write := func(at time.Time) error {
		_, err := fmt.Fprintf(output, "\r\x1b[2K%s", renderer.Text(terminalui.ToneWarning, fmt.Sprintf("● %s · elapsed %s", status, at.Sub(started).Round(time.Second))))
		return err
	}
	if err := write(started); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			_, err := fmt.Fprint(output, "\r\x1b[2K")
			if err != nil {
				return errors.Join(ctx.Err(), err)
			}
			return ctx.Err()
		case at := <-ticks:
			if err := write(at); err != nil {
				return err
			}
		}
	}
}

func clearCredentials(credentials *policy.Credentials) {
	clear(credentials.BasicPassword)
	clear(credentials.Password)
	*credentials = policy.Credentials{}
}

func isHostedPathPreflightError(err error) bool {
	if isPathPreflightError(err) || errors.Is(err, webdelivery.ErrEndpointPath) {
		return true
	}
	var remote *ipc.RemoteError
	return errors.As(err, &remote) && remote.Code == ipc.CodeInvalid && strings.Contains(remote.Message, webdelivery.ErrEndpointPath.Error())
}

func deliveryPolicy(plan operation.Plan) (delivery.Policy, error) {
	configured := delivery.DefaultPolicy()
	configured.Auth = delivery.AuthMode(plan.Options.Auth)
	configured.AuthAttempts = plan.Options.AuthAttempts
	configured.AuthFailAction = delivery.AuthFailAction(plan.Options.AuthFailAction)
	if !plan.Options.Limit.Unlimited {
		if plan.Options.Limit.Value > math.MaxInt64 {
			return delivery.Policy{}, errors.New("delivery limit exceeds the supported range")
		}
		configured.DeliveryLimit = delivery.Limit{Value: int64(plan.Options.Limit.Value)}
	}
	configured.AllowIP = make([]string, len(plan.Options.AllowIP))
	for index, prefix := range plan.Options.AllowIP {
		configured.AllowIP[index] = prefix.String()
	}
	configured.MaxFileSize = delivery.Limit{Unlimited: plan.Options.MaxFileSize.Unlimited, Value: plan.Options.MaxFileSize.Value}
	configured.MaxExtractedSize = delivery.Limit{Unlimited: plan.Options.MaxExtractedSize.Unlimited, Value: plan.Options.MaxExtractedSize.Value}
	configured.UploadRate = delivery.Limit{Unlimited: plan.Options.UploadRate.Unlimited, Value: plan.Options.UploadRate.Value}
	configured.DownloadRate = delivery.Limit{Unlimited: plan.Options.DownloadRate.Unlimited, Value: plan.Options.DownloadRate.Value}
	configured.NoUI = plan.Options.NoUI
	return configured, configured.Validate()
}
