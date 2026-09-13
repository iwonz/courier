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
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/webdelivery"
	"github.com/iwonz/courier/internal/worker"
	"github.com/pkg/sftp"
)

var (
	defaultStateDirectory = delivery.DefaultStateDirectory
	openDeliveryStore     = delivery.OpenStore
	newWebDefinition      = webdelivery.NewDefinition
	webAddress            = webdelivery.URL
	webNow                = time.Now
	newWebCoordinator     = func(store *delivery.Store, directory string) webCoordinator {
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
	}
	return newWebCoordinator(store, stateDirectory).Acquire(ctx, worker.AcquireRequest{
		Bind: plan.Options.Listen, Compatibility: "web-v1/" + buildinfo.Version, Route: route,
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
			if _, err := fmt.Fprint(output, "Basic username: "); err != nil {
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

func endpointCredentialProvider(factory sshx.Factory) func(context.Context, endpoint.Endpoint) (webdelivery.EndpointRuntime, error) {
	return func(ctx context.Context, value endpoint.Endpoint) (runtime webdelivery.EndpointRuntime, resultErr error) {
		if !value.Remote {
			return webdelivery.EndpointRuntime{}, nil
		}
		defer func() {
			if resultErr != nil {
				runtime.Clear()
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
			return runtime, err
		}
		if _, _, err := detectSSHPlatform(ctx, connection); err != nil {
			_ = closeSSHConnection(connection)
			return runtime, err
		}
		if err := closeSSHConnection(connection); err != nil {
			return runtime, err
		}
		return runtime, nil
	}
}

func webRunner(credentials func(context.Context, operation.AuthMode) (policy.Credentials, error), endpointCredentials func(context.Context, endpoint.Endpoint) (webdelivery.EndpointRuntime, error)) func(context.Context, operation.Plan, io.Writer) error {
	return func(ctx context.Context, plan operation.Plan, output io.Writer) error {
		if credentials == nil {
			return errors.New("web credential provider is unavailable")
		}
		remote := plan.Source
		if !remote.Remote {
			remote = plan.Destination
		}
		endpointRuntime := webdelivery.EndpointRuntime{}
		if remote.Remote {
			if endpointCredentials == nil {
				return errors.New("web SSH credential provider is unavailable")
			}
			var err error
			endpointRuntime, err = endpointCredentials(ctx, remote)
			if err != nil {
				return err
			}
		}
		defer endpointRuntime.Clear()
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
		address, err := webAddress(plan.Options.Listen, definition.Token)
		if err != nil {
			if acquired.Lease != nil {
				_ = releaseWebLease(acquired.Lease, context.Background())
			}
			return err
		}
		if _, err := fmt.Fprintf(output, "delivery: %s\nid: %s\n", address, acquired.DeliveryID); err != nil {
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
		<-ctx.Done()
		releaseContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return errors.Join(ctx.Err(), releaseWebLease(acquired.Lease, releaseContext))
	}
}

func clearCredentials(credentials *policy.Credentials) {
	clear(credentials.BasicPassword)
	clear(credentials.Password)
	*credentials = policy.Credentials{}
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
