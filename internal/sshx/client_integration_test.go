package sshx

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestNativeSSHAndSFTP(t *testing.T) {
	server := newTestSSHServer(t, true, true)
	factory := testFactory(t, server, "")
	connection, err := factory.Open(context.Background(), "target", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = connection.Close() })
	output, err := connection.Run(context.Background(), posixPlatformProbe)
	if err != nil || string(output) != "Linux\nx86_64\n" {
		t.Fatalf("output=%q err=%v", output, err)
	}
	backend := NewSFTPBackend(connection.SFTP)
	if err := backend.MkdirAll("/dir", 0o700); err != nil {
		t.Fatal(err)
	}
	writer, err := backend.Create("/dir/file", 0o640)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("remote data")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Sync(); err == nil {
		t.Fatal("expected sync error after close")
	}
	if err := backend.Chmod("/dir/file", 0o600); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(1_700_000_000, 0)
	if err := backend.Chtimes("/dir/file", stamp, stamp); err != nil {
		t.Fatal(err)
	}
	if info, err := backend.Lstat("/dir/file"); err != nil || info.Size() != int64(len("remote data")) {
		t.Fatalf("info=%v err=%v", info, err)
	}
	entries, err := backend.ReadDir("/dir")
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	reader, err := backend.Open("/dir/file")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(reader)
	_ = reader.Close()
	if string(data) != "remote data" {
		t.Fatalf("data=%q", data)
	}
	if err := backend.Symlink("/dir/file", "/link"); err != nil {
		t.Fatal(err)
	}
	if target, err := backend.Readlink("/link"); err != nil || target != "/dir/file" {
		t.Fatalf("target=%q err=%v", target, err)
	}
	if err := backend.Rename("/dir/file", "/dir/renamed"); err != nil {
		t.Fatal(err)
	}
	if err := (realSFTP{client: connection.SFTP}).Rename("/dir/renamed", "/dir/file"); err != nil {
		t.Fatal(err)
	}
	if err := backend.RemoveAll("/dir"); err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.Create("/closed", 0o600); err == nil {
		t.Fatal("expected create error after close")
	}
}

func TestPasswordFallbackAndHostVerification(t *testing.T) {
	server := newTestSSHServer(t, false, true)
	factory := testFactory(t, server, "")
	factory.ReadFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
	secret := []byte("secret")
	factory.Prompt = func(string) ([]byte, error) { return secret, nil }
	connection, err := factory.Open(context.Background(), "target", "")
	if err != nil {
		t.Fatal(err)
	}
	_ = connection.Close()
	if string(secret) != "\x00\x00\x00\x00\x00\x00" {
		t.Fatalf("password was not zeroed: %v", secret)
	}

	wrongKnownHosts := filepath.Join(t.TempDir(), "known_hosts")
	_, otherPrivate, _ := ed25519.GenerateKey(rand.Reader)
	otherSigner, _ := ssh.NewSignerFromKey(otherPrivate)
	line := knownhosts.Line([]string{knownhosts.Normalize(server.address)}, otherSigner.PublicKey())
	if err := os.WriteFile(wrongKnownHosts, []byte(line+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	factory = testFactory(t, server, wrongKnownHosts)
	if _, err := factory.Open(context.Background(), "target", ""); err == nil || strings.Contains(err.Error(), "secret") {
		t.Fatalf("expected safe host-key failure, got %v", err)
	}
}

func TestProxyJump(t *testing.T) {
	bastion := newTestSSHServer(t, true, true)
	target := newTestSSHServer(t, true, true)
	root := t.TempDir()
	identity := writeIdentity(t, root, target.clientPrivate, nil)
	known := filepath.Join(root, "known_hosts")
	knownData := bastion.knownHostsLine + "\n" + target.knownHostsLine + "\n"
	if err := os.WriteFile(known, []byte(knownData), 0o600); err != nil {
		t.Fatal(err)
	}
	bastionHost, bastionPort, _ := net.SplitHostPort(bastion.address)
	targetHost, targetPort, _ := net.SplitHostPort(target.address)
	configPath := filepath.Join(root, "config")
	configText := fmt.Sprintf("Host target\n HostName %s\n Port %s\n User test\n IdentityFile %s\n ProxyJump bastion:%s\nHost bastion\n HostName %s\n Port %s\n User test\n IdentityFile %s\n", targetHost, targetPort, identity, bastionPort, bastionHost, bastionPort, identity)
	if err := os.WriteFile(configPath, []byte(configText), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(configPath, root)
	if err != nil {
		t.Fatal(err)
	}
	connection, err := (Factory{Config: config, DefaultUser: "test", KnownHosts: known}).Open(context.Background(), "target", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(connection.hops) != 2 {
		t.Fatalf("hops=%d", len(connection.hops))
	}
	_ = connection.Close()
}

func TestConnectionRunFailures(t *testing.T) {
	empty := &Connection{}
	if _, err := empty.Run(context.Background(), "command"); err == nil {
		t.Fatal("expected closed error")
	}
	server := newTestSSHServer(t, true, true)
	connection, err := testFactory(t, server, "").Open(context.Background(), "target", "")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := connection.Run(ctx, "block"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	_ = connection.Close()
	if _, err := connection.Run(context.Background(), "command"); err == nil {
		t.Fatal("expected session error")
	}
}

func TestIdentityAndAgentAuthMethods(t *testing.T) {
	_, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	root := t.TempDir()
	encrypted := writeIdentity(t, root, privateKey, []byte("passphrase"))
	secret := []byte("passphrase")
	factory := Factory{Prompt: func(string) ([]byte, error) { return secret, nil }}
	methods, closers, err := factory.authMethods(Target{IdentityFiles: []string{encrypted, filepath.Join(root, "missing")}})
	if err != nil || len(methods) != 1 || len(closers) != 0 {
		t.Fatalf("methods=%d closers=%d err=%v", len(methods), len(closers), err)
	}
	for _, value := range secret {
		if value != 0 {
			t.Fatal("passphrase was not zeroed")
		}
	}
	serverSide, clientSide := net.Pipe()
	keyring := agent.NewKeyring()
	if err := keyring.Add(agent.AddedKey{PrivateKey: privateKey}); err != nil {
		t.Fatal(err)
	}
	go func() { _ = agent.ServeAgent(keyring, serverSide) }()
	factory = Factory{AgentSocket: "agent", DialAgent: func(string) (net.Conn, error) { return clientSide, nil }}
	methods, closers, err = factory.authMethods(Target{})
	if err != nil || len(methods) != 1 || len(closers) != 1 {
		t.Fatalf("agent methods=%d closers=%d err=%v", len(methods), len(closers), err)
	}
	if err := (&Connection{agents: closers}).Close(); err != nil {
		t.Fatal(err)
	}
	_ = serverSide.Close()
}

func TestClientFailurePaths(t *testing.T) {
	cause := errors.New("capability")
	capability := &CapabilityError{Capability: "sftp", Cause: cause}
	if !strings.Contains(capability.Error(), "sftp") || !errors.Is(capability, cause) {
		t.Fatalf("capability error=%v", capability)
	}
	if _, err := (Factory{}).Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected missing config")
	}
	config := EmptyConfig(t.TempDir())
	if _, err := (Factory{Config: config}).Open(context.Background(), "", ""); err == nil {
		t.Fatal("expected resolve error")
	}
	badJumpRoot := t.TempDir()
	badJumpConfig := configFromText(t, badJumpRoot, "Host target\n ProxyJump jump\nHost jump\n Port bad\n")
	if _, err := (Factory{Config: badJumpConfig, DefaultUser: "test", KnownHosts: "known"}).Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected jump resolve error")
	}
	jumpRoot := t.TempDir()
	jumpFailureConfig := configFromText(t, jumpRoot, "Host target\n ProxyJump jump\n")
	if _, err := (Factory{Config: jumpFailureConfig, DefaultUser: "test", KnownHosts: filepath.Join(t.TempDir(), "missing")}).Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected jump connection error")
	}
	readFailure := Factory{ReadFile: func(string) ([]byte, error) { return nil, errors.New("read identity") }}
	if _, _, err := readFailure.authMethods(Target{IdentityFiles: []string{"identity"}}); err == nil {
		t.Fatal("expected identity read error")
	}
	server := newTestSSHServer(t, true, false)
	readFailure = testFactory(t, server, "")
	readFailure.ReadFile = func(string) ([]byte, error) { return nil, errors.New("read identity") }
	if _, err := readFailure.Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected connection identity error")
	}
	invalidKey := Factory{ReadFile: func(string) ([]byte, error) { return []byte("invalid"), nil }}
	if _, _, err := invalidKey.authMethods(Target{IdentityFiles: []string{"identity"}}); err == nil {
		t.Fatal("expected identity parse error")
	}
	_, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	encrypted := writeIdentity(t, t.TempDir(), privateKey, []byte("pass"))
	promptFailure := Factory{Prompt: func(string) ([]byte, error) { return nil, errors.New("prompt") }}
	if _, _, err := promptFailure.authMethods(Target{IdentityFiles: []string{encrypted}}); err == nil {
		t.Fatal("expected passphrase prompt error")
	}
	defaultAgent := Factory{AgentSocket: filepath.Join(t.TempDir(), "missing-agent")}
	if methods, closers, err := defaultAgent.authMethods(Target{}); err != nil || methods != nil || len(closers) != 0 {
		t.Fatalf("default agent result=%v,%v,%v", methods, closers, err)
	}
	serverSide, clientSide := net.Pipe()
	_ = serverSide.Close()
	brokenAgent := Factory{AgentSocket: "agent", DialAgent: func(string) (net.Conn, error) { return clientSide, nil }}
	if methods, closers, err := brokenAgent.authMethods(Target{}); err != nil || methods != nil || len(closers) != 0 {
		t.Fatalf("broken agent result=%v,%v,%v", methods, closers, err)
	}
	target := Target{Host: "127.0.0.1", Port: 1, User: "test", KnownHostsFile: "known"}
	configFailure := Factory{ClientConfig: func(Target, []ssh.AuthMethod) (*ssh.ClientConfig, error) { return nil, errors.New("config") }}
	if _, err := configFailure.dialSSH(context.Background(), target, nil, nil); err == nil {
		t.Fatal("expected client config error")
	}
	dialFailure := Factory{
		ClientConfig: func(Target, []ssh.AuthMethod) (*ssh.ClientConfig, error) { return &ssh.ClientConfig{}, nil },
		DialContext:  func(context.Context, string, string) (net.Conn, error) { return nil, errors.New("dial") },
	}
	if _, err := dialFailure.dialSSH(context.Background(), target, nil, nil); err == nil {
		t.Fatal("expected network dial error")
	}
	if _, err := secureClientConfig(Target{}, nil); err == nil {
		t.Fatal("expected missing user")
	}
	if _, err := secureClientConfig(Target{User: "test"}, nil); err == nil {
		t.Fatal("expected missing known_hosts")
	}
	if _, err := secureClientConfig(Target{User: "test", KnownHostsFile: filepath.Join(t.TempDir(), "missing")}, nil); err == nil {
		t.Fatal("expected known_hosts load error")
	}
}

func TestPasswordPromptFailureAndSFTPCapability(t *testing.T) {
	passwordServer := newTestSSHServer(t, false, true)
	factory := testFactory(t, passwordServer, "")
	factory.ReadFile = func(string) ([]byte, error) { return nil, os.ErrNotExist }
	factory.Prompt = func(string) ([]byte, error) { return nil, errors.New("prompt") }
	if _, err := factory.Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected password prompt error")
	}
	server := newTestSSHServer(t, true, false)
	factory = testFactory(t, server, "")
	factory.NewSFTP = func(*ssh.Client) (*sftp.Client, error) { return nil, errors.New("no sftp") }
	if _, err := factory.Open(context.Background(), "target", ""); err == nil {
		t.Fatal("expected SFTP capability error")
	} else {
		var capability *CapabilityError
		if !errors.As(err, &capability) {
			t.Fatalf("error=%v", err)
		}
	}
}

func testFactory(t *testing.T, server *testSSHServer, knownOverride string) Factory {
	t.Helper()
	root := t.TempDir()
	identity := writeIdentity(t, root, server.clientPrivate, nil)
	known := knownOverride
	if known == "" {
		known = filepath.Join(root, "known_hosts")
		if err := os.WriteFile(known, []byte(server.knownHostsLine+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	host, port, _ := net.SplitHostPort(server.address)
	configPath := filepath.Join(root, "config")
	configText := fmt.Sprintf("Host target\n HostName %s\n Port %s\n User test\n IdentityFile %s\n", host, port, identity)
	if err := os.WriteFile(configPath, []byte(configText), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := LoadConfig(configPath, root)
	if err != nil {
		t.Fatal(err)
	}
	return Factory{Config: config, DefaultUser: "test", KnownHosts: known}
}

func writeIdentity(t *testing.T, root string, privateKey ed25519.PrivateKey, passphrase []byte) string {
	t.Helper()
	var block *pem.Block
	var err error
	if passphrase == nil {
		block, err = ssh.MarshalPrivateKey(privateKey, "courier test")
	} else {
		block, err = ssh.MarshalPrivateKeyWithPassphrase(privateKey, "courier test", passphrase)
	}
	if err != nil {
		t.Fatal(err)
	}
	name := filepath.Join(root, "identity")
	if err := os.WriteFile(name, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	return name
}

type testSSHServer struct {
	listener       net.Listener
	address        string
	knownHostsLine string
	clientPrivate  ed25519.PrivateKey
	configuration  *ssh.ServerConfig
	handlers       sftp.Handlers
	done           chan struct{}
	once           sync.Once
}

func newTestSSHServer(t *testing.T, allowPublic, allowPassword bool) *testSSHServer {
	t.Helper()
	_, hostPrivate, _ := ed25519.GenerateKey(rand.Reader)
	hostSigner, err := ssh.NewSignerFromKey(hostPrivate)
	if err != nil {
		t.Fatal(err)
	}
	_, clientPrivate, _ := ed25519.GenerateKey(rand.Reader)
	configuration := &ssh.ServerConfig{}
	if allowPublic {
		configuration.PublicKeyCallback = func(_ ssh.ConnMetadata, _ ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		}
	}
	if allowPassword {
		configuration.PasswordCallback = func(_ ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if string(password) == "secret" {
				return nil, nil
			}
			return nil, errors.New("wrong password")
		}
	}
	configuration.AddHostKey(hostSigner)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &testSSHServer{
		listener: listener, address: listener.Addr().String(), clientPrivate: clientPrivate,
		configuration: configuration, handlers: sftp.InMemHandler(), done: make(chan struct{}),
	}
	server.handlers.FileCmd = testFileCommand{base: server.handlers.FileCmd}
	server.knownHostsLine = knownhosts.Line([]string{knownhosts.Normalize(server.address)}, hostSigner.PublicKey())
	go server.serve()
	t.Cleanup(server.close)
	return server
}

type testFileCommand struct{ base sftp.FileCmder }

func (t testFileCommand) Filecmd(request *sftp.Request) error {
	if request.Method == "Setstat" {
		return nil
	}
	return t.base.Filecmd(request)
}

func (t testFileCommand) PosixRename(request *sftp.Request) error {
	return t.base.(sftp.PosixRenameFileCmder).PosixRename(request)
}

func (s *testSSHServer) close() {
	s.once.Do(func() {
		_ = s.listener.Close()
		<-s.done
	})
}

func (s *testSSHServer) serve() {
	defer close(s.done)
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(connection)
	}
}

func (s *testSSHServer) handleConnection(connection net.Conn) {
	serverConnection, channels, requests, err := ssh.NewServerConn(connection, s.configuration)
	if err != nil {
		_ = connection.Close()
		return
	}
	defer serverConnection.Close()
	go ssh.DiscardRequests(requests)
	for channelRequest := range channels {
		switch channelRequest.ChannelType() {
		case "session":
			channel, requests, err := channelRequest.Accept()
			if err == nil {
				go s.handleSession(channel, requests)
			}
		case "direct-tcpip":
			go s.handleDirect(channelRequest)
		default:
			_ = channelRequest.Reject(ssh.UnknownChannelType, "unsupported")
		}
	}
}

func (s *testSSHServer) handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	for request := range requests {
		switch request.Type {
		case "subsystem":
			var payload struct{ Name string }
			_ = ssh.Unmarshal(request.Payload, &payload)
			if payload.Name != "sftp" {
				_ = request.Reply(false, nil)
				continue
			}
			_ = request.Reply(true, nil)
			server := sftp.NewRequestServer(channel, s.handlers)
			_ = server.Serve()
			_ = server.Close()
			return
		case "exec":
			var payload struct{ Command string }
			_ = ssh.Unmarshal(request.Payload, &payload)
			_ = request.Reply(true, nil)
			if payload.Command == "block" {
				continue
			}
			if payload.Command == posixPlatformProbe {
				_, _ = io.WriteString(channel, "Linux\nx86_64\n")
			}
			_, _ = channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
			_ = channel.Close()
			return
		default:
			_ = request.Reply(false, nil)
		}
	}
}

func (s *testSSHServer) handleDirect(request ssh.NewChannel) {
	var payload struct {
		Host       string
		Port       uint32
		OriginHost string
		OriginPort uint32
	}
	if err := ssh.Unmarshal(request.ExtraData(), &payload); err != nil {
		_ = request.Reject(ssh.ConnectionFailed, "invalid target")
		return
	}
	destination, err := net.Dial("tcp", net.JoinHostPort(payload.Host, fmt.Sprint(payload.Port)))
	if err != nil {
		_ = request.Reject(ssh.ConnectionFailed, err.Error())
		return
	}
	channel, requests, err := request.Accept()
	if err != nil {
		_ = destination.Close()
		return
	}
	go ssh.DiscardRequests(requests)
	go func() { _, _ = io.Copy(channel, destination); _ = channel.CloseWrite() }()
	_, _ = io.Copy(destination, channel)
	_ = destination.Close()
	_ = channel.Close()
}
