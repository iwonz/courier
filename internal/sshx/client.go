package sshx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SecretPrompt obtains a password or key passphrase without command-line use.
type SecretPrompt func(label string) ([]byte, error)

// Factory resolves configuration and opens independently authenticated clients.
type Factory struct {
	Config       *Config
	DefaultUser  string
	KnownHosts   string
	AgentSocket  string
	Prompt       SecretPrompt
	ReadFile     func(string) ([]byte, error)
	DialContext  func(context.Context, string, string) (net.Conn, error)
	DialAgent    func(string) (net.Conn, error)
	NewSFTP      func(*ssh.Client) (*sftp.Client, error)
	SFTPFallback func(context.Context, *Connection, *CapabilityError) (*sftp.Client, io.Closer, error)
	ClientConfig func(Target, []ssh.AuthMethod) (*ssh.ClientConfig, error)
}

// CapabilityError identifies a remote feature gap eligible for an optional helper.
type CapabilityError struct {
	Capability string
	Cause      error
}

func (e *CapabilityError) Error() string {
	return fmt.Sprintf("remote capability %s is unavailable: %v", e.Capability, e.Cause)
}
func (e *CapabilityError) Unwrap() error { return e.Cause }

// Connection owns SFTP and every SSH hop used to reach it.
type Connection struct {
	Target               Target
	SFTP                 *sftp.Client
	SSH                  *ssh.Client
	hops                 []*ssh.Client
	agents               []io.Closer
	helpers              []io.Closer
	helperSessionFactory func() (helperSession, error)
	helperCommandRunner  func(context.Context, string) ([]byte, error)
}

// Close releases SFTP, SSH hops in reverse order, and agent connections.
func (c *Connection) Close() error {
	var result error
	if c.SFTP != nil {
		result = errors.Join(result, c.SFTP.Close())
	}
	for index := len(c.helpers) - 1; index >= 0; index-- {
		result = errors.Join(result, c.helpers[index].Close())
	}
	for index := len(c.hops) - 1; index >= 0; index-- {
		result = errors.Join(result, c.hops[index].Close())
	}
	for _, closer := range c.agents {
		result = errors.Join(result, closer.Close())
	}
	return result
}

// Run executes a fixed remote capability probe with cancellation.
func (c *Connection) Run(ctx context.Context, command string) ([]byte, error) {
	if c != nil && c.helperCommandRunner != nil {
		return c.helperCommandRunner(ctx, command)
	}
	if c == nil || c.SSH == nil {
		return nil, errors.New("SSH connection is closed")
	}
	session, err := c.SSH.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	type outcome struct {
		data []byte
		err  error
	}
	done := make(chan outcome, 1)
	go func() {
		data, err := session.CombinedOutput(command)
		done <- outcome{data: data, err: err}
	}()
	select {
	case <-ctx.Done():
		_ = session.Close()
		return nil, ctx.Err()
	case result := <-done:
		return result.data, result.err
	}
}

// Open creates a verified ProxyJump chain and SFTP client.
func (f Factory) Open(ctx context.Context, alias, explicitUser string) (*Connection, error) {
	if f.Config == nil {
		return nil, errors.New("SSH config is required")
	}
	target, err := f.Config.Resolve(alias, explicitUser, f.DefaultUser, f.KnownHosts)
	if err != nil {
		return nil, err
	}
	connection := &Connection{Target: target}
	var upstream *ssh.Client
	for _, jump := range target.ProxyJump {
		jumpTarget, err := f.Config.Resolve(jump.Host, jump.User, f.DefaultUser, f.KnownHosts)
		if err != nil {
			_ = connection.Close()
			return nil, err
		}
		if jump.Port != 0 {
			jumpTarget.Port = jump.Port
		}
		client, agentClosers, err := f.connect(ctx, jumpTarget, upstream)
		connection.agents = append(connection.agents, agentClosers...)
		if err != nil {
			_ = connection.Close()
			return nil, fmt.Errorf("connect jump %s: %w", jumpTarget.Alias, err)
		}
		connection.hops = append(connection.hops, client)
		upstream = client
	}
	client, agentClosers, err := f.connect(ctx, target, upstream)
	connection.agents = append(connection.agents, agentClosers...)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	connection.SSH = client
	connection.hops = append(connection.hops, client)
	newSFTP := f.NewSFTP
	if newSFTP == nil {
		newSFTP = func(client *ssh.Client) (*sftp.Client, error) { return sftp.NewClient(client) }
	}
	connection.SFTP, err = newSFTP(client)
	if err != nil {
		capability := &CapabilityError{Capability: "sftp", Cause: err}
		if f.SFTPFallback == nil {
			_ = connection.Close()
			return nil, capability
		}
		var runtime io.Closer
		connection.SFTP, runtime, err = f.SFTPFallback(ctx, connection, capability)
		if err != nil {
			_ = connection.Close()
			return nil, err
		}
		if connection.SFTP == nil || runtime == nil {
			_ = connection.Close()
			return nil, errors.New("SFTP fallback returned incomplete runtime")
		}
		connection.helpers = append(connection.helpers, runtime)
	}
	return connection, nil
}

func (f Factory) connect(ctx context.Context, target Target, upstream *ssh.Client) (*ssh.Client, []io.Closer, error) {
	authMethods, closers, err := f.authMethods(target)
	if err != nil {
		return nil, closers, err
	}
	client, err := f.dialSSH(ctx, target, authMethods, upstream)
	if err == nil || f.Prompt == nil || !isAuthenticationError(err) {
		return client, closers, err
	}
	secret, promptErr := f.Prompt(fmt.Sprintf("Password for %s@%s: ", target.User, target.Alias))
	if promptErr != nil {
		return nil, closers, promptErr
	}
	client, err = f.dialSSH(ctx, target, append(authMethods, ssh.Password(string(secret))), upstream)
	for index := range secret {
		secret[index] = 0
	}
	return client, closers, err
}

func (f Factory) authMethods(target Target) ([]ssh.AuthMethod, []io.Closer, error) {
	readFile := f.ReadFile
	if readFile == nil {
		readFile = os.ReadFile
	}
	var signers []ssh.Signer
	for _, identity := range target.IdentityFiles {
		data, err := readFile(identity)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		signer, err := ssh.ParsePrivateKey(data)
		var missing *ssh.PassphraseMissingError
		if errors.As(err, &missing) && f.Prompt != nil {
			secret, promptErr := f.Prompt(fmt.Sprintf("Passphrase for %s: ", identity))
			if promptErr != nil {
				return nil, nil, promptErr
			}
			signer, err = ssh.ParsePrivateKeyWithPassphrase(data, secret)
			for index := range secret {
				secret[index] = 0
			}
		}
		if err != nil {
			return nil, nil, fmt.Errorf("parse identity %q: %w", identity, err)
		}
		signers = append(signers, signer)
	}
	var closers []io.Closer
	if f.AgentSocket != "" {
		dialAgent := f.DialAgent
		if dialAgent == nil {
			dialAgent = defaultAgentDial
		}
		connection, err := dialAgent(f.AgentSocket)
		if err == nil {
			agentSigners, signerErr := agent.NewClient(connection).Signers()
			if signerErr == nil {
				closers = append(closers, connection)
				signers = append(signers, agentSigners...)
			} else {
				_ = connection.Close()
			}
		}
	}
	if len(signers) == 0 {
		return nil, closers, nil
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signers...)}, closers, nil
}

func (f Factory) dialSSH(ctx context.Context, target Target, authentication []ssh.AuthMethod, upstream *ssh.Client) (*ssh.Client, error) {
	clientConfig := f.ClientConfig
	if clientConfig == nil {
		clientConfig = secureClientConfig
	}
	configuration, err := clientConfig(target, authentication)
	if err != nil {
		return nil, err
	}
	var connection net.Conn
	if upstream == nil {
		dial := f.DialContext
		if dial == nil {
			dialer := &net.Dialer{Timeout: 15 * time.Second}
			dial = func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, address)
			}
		}
		connection, err = dial(ctx, "tcp", target.Address())
	} else {
		connection, err = upstream.DialContext(ctx, "tcp", target.Address())
	}
	if err != nil {
		return nil, err
	}
	channel, channels, requests, err := ssh.NewClientConn(connection, target.Address(), configuration)
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	return ssh.NewClient(channel, channels, requests), nil
}

func secureClientConfig(target Target, authentication []ssh.AuthMethod) (*ssh.ClientConfig, error) {
	if target.User == "" {
		return nil, errors.New("SSH user is required")
	}
	if target.KnownHostsFile == "" {
		return nil, errors.New("known_hosts file is required")
	}
	callback, err := knownhosts.New(target.KnownHostsFile)
	if err != nil {
		return nil, fmt.Errorf("load known_hosts: %w", err)
	}
	return &ssh.ClientConfig{User: target.User, Auth: authentication, HostKeyCallback: callback, Timeout: 15 * time.Second}, nil
}

func isAuthenticationError(err error) bool {
	var authError *ssh.ServerAuthError
	return errors.As(err, &authError) || strings.Contains(strings.ToLower(err.Error()), "unable to authenticate")
}
