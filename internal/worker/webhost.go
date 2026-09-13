package worker

import (
	"context"
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/iwonz/courier/internal/buildinfo"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/helper"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/selection"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/webdelivery"
)

var (
	workerHome       = os.UserHomeDir
	workerUser       = user.Current
	workerLoadConfig = sshx.LoadConfig
	workerOpenSSH    = func(ctx context.Context, factory sshx.Factory, value endpoint.Endpoint) (*sshx.Connection, error) {
		return factory.Open(ctx, value.Host, value.User)
	}
	workerDetectSSH      = sshx.DetectPlatform
	workerOpenRooted     = fsx.OpenRootedForPath
	workerRemoteResource = func(connection *sshx.Connection, value endpoint.Endpoint) *webdelivery.Resource {
		value.Host = connection.Target.Host
		value.User = connection.Target.User
		return &webdelivery.Resource{Endpoint: value, Backend: sshx.NewSFTPBackend(connection.SFTP), Path: value.Path, Close: connection.Close}
	}
)

func newWorkerWebHost(stop func(delivery.ID)) (DeliveryHost, error) {
	home, err := workerHome()
	if err != nil {
		return nil, err
	}
	configuration, err := workerLoadConfig(filepath.Join(home, ".ssh", "config"), home)
	if errors.Is(err, os.ErrNotExist) {
		configuration = sshx.EmptyConfig(home)
	} else if err != nil {
		return nil, err
	}
	username := ""
	if current, currentErr := workerUser(); currentErr == nil {
		username = current.Username
	}
	factory := sshx.Factory{
		Config: configuration, DefaultUser: username, KnownHosts: filepath.Join(home, ".ssh", "known_hosts"),
		AgentSocket: sshx.AgentEndpoint(os.Getenv("SSH_AUTH_SOCK")),
	}
	helperSource := helper.Source{Repository: "iwonz/courier", Version: buildinfo.Version}
	helperManager := &helper.Manager{Confirm: approveWorkerHelper, Acquire: helperSource.Acquire, Deploy: sshx.DeployHelper}
	opener := func(ctx context.Context, value endpoint.Endpoint, endpointRuntime webdelivery.EndpointRuntime) (*webdelivery.Resource, error) {
		if !value.Remote {
			backend, relative, openErr := workerOpenRooted(value.Path)
			if openErr != nil {
				return nil, openErr
			}
			return &webdelivery.Resource{Endpoint: value, Backend: backend, Path: relative, Close: backend.Close}, nil
		}
		connectionFactory := factory
		if len(endpointRuntime.Credentials) != 0 {
			connectionFactory.Prompt = endpointRuntime.Prompt
		}
		if endpointRuntime.AllowHelper {
			connectionFactory.SFTPFallback = helperManager.Fallback
		}
		connection, openErr := workerOpenSSH(ctx, connectionFactory, value)
		if openErr != nil {
			return nil, openErr
		}
		if _, _, detectErr := workerDetectSSH(ctx, connection); detectErr != nil {
			_ = connection.Close()
			return nil, detectErr
		}
		return workerRemoteResource(connection, value), nil
	}
	return webdelivery.NewHost(webdelivery.HostOptions{
		Open: opener,
		Select: func(rules []operation.SelectionRule) (selection.Selector, error) {
			return selection.Compile(rules, selection.OpenFile)
		},
		PolicyDependencies: policy.DefaultDependencies(),
		StopRequested: func(id delivery.ID) {
			if stop != nil {
				stop(id)
			}
		},
	})
}

func approveWorkerHelper(context.Context, string) (bool, error) { return true, nil }
