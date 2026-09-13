package sshx

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pkg/sftp"
)

func TestHelperCommandsAndPaths(t *testing.T) {
	originalRandom := randomHelperBytes
	t.Cleanup(func() { randomHelperBytes = originalRandom })
	randomHelperBytes = func(data []byte) (int, error) {
		for index := range data {
			data[index] = 0xab
		}
		return len(data), nil
	}
	posixPath, err := newRemoteHelperPath(Platform{OS: "linux"})
	if err != nil || posixPath != "/tmp/courier-helper-abababababababababababababababab/courier" {
		t.Fatalf("POSIX path=%q err=%v", posixPath, err)
	}
	windowsPath, err := newRemoteHelperPath(Platform{OS: "windows"})
	if err != nil || windowsPath != `courier-helper-abababababababababababababababab\courier.exe` {
		t.Fatalf("Windows path=%q err=%v", windowsPath, err)
	}
	digest := strings.Repeat("0", 64)
	for _, command := range []string{helperUploadCommand(Platform{OS: "linux"}, posixPath, digest), helperStartCommand(Platform{OS: "linux"}, posixPath), helperCleanupCommand(Platform{OS: "linux"}, posixPath)} {
		if !strings.Contains(command, "courier-helper-") {
			t.Fatalf("missing helper path in %q", command)
		}
	}
	for _, command := range []string{helperUploadCommand(Platform{OS: "windows"}, windowsPath, digest), helperStartCommand(Platform{OS: "windows"}, windowsPath), helperCleanupCommand(Platform{OS: "windows"}, windowsPath)} {
		if !strings.HasPrefix(command, "powershell -NoProfile -NonInteractive -EncodedCommand ") {
			t.Fatalf("invalid PowerShell command %q", command)
		}
	}
	if shellQuote("a'b") != `'a'\''b'` || powerShellQuote("a'b") != `'a''b'` {
		t.Fatal("command quoting failed")
	}
	encoded := encodedPowerShell("A😀")
	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encoded, "powershell -NoProfile -NonInteractive -EncodedCommand "))
	if err != nil || len(payload) != 6 {
		t.Fatalf("encoded payload length=%d err=%v", len(payload), err)
	}
	randomHelperBytes = func([]byte) (int, error) { return 0, errors.New("random") }
	if _, err := newRemoteHelperPath(Platform{}); err == nil {
		t.Fatal("expected random error")
	}
}

func TestDeployHelperFailures(t *testing.T) {
	isolateHelperHooks(t)
	digest := strings.Repeat("0", 64)
	if _, _, err := DeployHelper(context.Background(), &Connection{}, Platform{OS: "linux"}, "binary", "bad"); err == nil {
		t.Fatal("expected digest error")
	}
	randomHelperBytes = func([]byte) (int, error) { return 0, errors.New("random") }
	if _, _, err := DeployHelper(context.Background(), &Connection{}, Platform{OS: "linux"}, "binary", digest); err == nil {
		t.Fatal("expected random helper path error")
	}
	randomHelperBytes = func(data []byte) (int, error) { return len(data), nil }
	if _, _, err := DeployHelper(context.Background(), nil, Platform{OS: "linux"}, "binary", digest); err == nil {
		t.Fatal("expected closed connection error")
	}
	connection := &Connection{helperSessionFactory: func() (helperSession, error) { return &fakeHelperSession{}, nil }}
	openHelperBinary = func(string) (io.ReadCloser, error) { return nil, errors.New("open") }
	if _, _, err := DeployHelper(context.Background(), connection, Platform{OS: "linux"}, "binary", digest); err == nil {
		t.Fatal("expected binary open error")
	}

	cleanups := 0
	connection.helperCommandRunner = func(context.Context, string) ([]byte, error) { cleanups++; return nil, nil }
	openHelperBinary = func(string) (io.ReadCloser, error) { return &testReadCloser{Reader: strings.NewReader("binary")}, nil }
	connection.helperSessionFactory = func() (helperSession, error) { return &fakeHelperSession{combinedErr: errors.New("upload")}, nil }
	if _, _, err := DeployHelper(context.Background(), connection, Platform{OS: "linux"}, "binary", digest); err == nil || cleanups != 1 {
		t.Fatalf("expected upload cleanup: count=%d err=%v", cleanups, err)
	}
	openHelperBinary = func(string) (io.ReadCloser, error) {
		return &testReadCloser{Reader: strings.NewReader("binary"), closeErr: errors.New("close")}, nil
	}
	connection.helperSessionFactory = func() (helperSession, error) { return &fakeHelperSession{}, nil }
	if _, _, err := DeployHelper(context.Background(), connection, Platform{OS: "linux"}, "binary", digest); err == nil || cleanups != 2 {
		t.Fatalf("expected source close cleanup: count=%d err=%v", cleanups, err)
	}

	blocked := make(chan struct{})
	session := &fakeHelperSession{combinedBlock: blocked}
	connection.helperSessionFactory = func() (helperSession, error) { return session, nil }
	openHelperBinary = func(string) (io.ReadCloser, error) { return &testReadCloser{Reader: strings.NewReader("binary")}, nil }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := DeployHelper(ctx, connection, Platform{OS: "linux"}, "binary", digest); !errors.Is(err, context.Canceled) || !session.closed {
		t.Fatalf("expected canceled upload: closed=%v err=%v", session.closed, err)
	}
}

func TestNewHelperSessionAfterSSHClose(t *testing.T) {
	server := newTestSSHServer(t, true, true)
	connection, err := testFactory(t, server, "").Open(context.Background(), "target", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := connection.newHelperSession(); err == nil {
		t.Fatal("expected session creation error after close")
	}
}

func TestStartHelperFailuresAndRuntime(t *testing.T) {
	isolateHelperHooks(t)
	platform := Platform{OS: "linux"}
	remotePath := "/tmp/courier-helper-00/courier"
	cleanupCount := 0
	connection := &Connection{helperCommandRunner: func(context.Context, string) ([]byte, error) { cleanupCount++; return nil, nil }}
	connection.helperSessionFactory = func() (helperSession, error) { return nil, errors.New("session") }
	if _, _, err := startHelper(context.Background(), connection, platform, remotePath); err == nil || cleanupCount != 1 {
		t.Fatalf("expected session failure: cleanup=%d err=%v", cleanupCount, err)
	}

	for _, test := range []struct {
		name    string
		session *fakeHelperSession
	}{
		{"stdin", &fakeHelperSession{stdinErr: errors.New("stdin")}},
		{"stdout", &fakeHelperSession{stdoutErr: errors.New("stdout")}},
		{"start", &fakeHelperSession{startErr: errors.New("start")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			connection.helperSessionFactory = func() (helperSession, error) { return test.session, nil }
			if _, _, err := startHelper(context.Background(), connection, platform, remotePath); err == nil || !test.session.closed {
				t.Fatalf("closed=%v err=%v", test.session.closed, err)
			}
		})
	}

	connection.helperSessionFactory = func() (helperSession, error) { return &fakeHelperSession{}, nil }
	newPipeSFTP = func(io.Reader, io.WriteCloser, ...sftp.ClientOption) (*sftp.Client, error) {
		return nil, errors.New("pipe")
	}
	if _, _, err := startHelper(context.Background(), connection, platform, remotePath); err == nil {
		t.Fatal("expected pipe error")
	}

	blockPipe := make(chan struct{})
	pipeDone := make(chan struct{})
	session := &fakeHelperSession{closeSignal: blockPipe}
	connection.helperSessionFactory = func() (helperSession, error) { return session, nil }
	newPipeSFTP = func(io.Reader, io.WriteCloser, ...sftp.ClientOption) (*sftp.Client, error) {
		<-blockPipe
		close(pipeDone)
		return nil, errors.New("closed")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := startHelper(ctx, connection, platform, remotePath); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation: %v", err)
	}
	<-pipeDone

	session = &fakeHelperSession{waitErr: errors.New("wait")}
	connection.helperSessionFactory = func() (helperSession, error) { return session, nil }
	client := &sftp.Client{}
	newPipeSFTP = func(io.Reader, io.WriteCloser, ...sftp.ClientOption) (*sftp.Client, error) { return client, nil }
	gotClient, closer, err := startHelper(context.Background(), connection, platform, remotePath)
	if err != nil || gotClient != client {
		t.Fatalf("client=%v err=%v", gotClient, err)
	}
	if err := closer.Close(); err == nil {
		t.Fatal("expected runtime wait error")
	}
	countAfterFirst := cleanupCount
	if err := closer.Close(); err == nil || cleanupCount != countAfterFirst {
		t.Fatal("runtime close should be idempotent")
	}
	if err := (&helperRuntime{}).Close(); err != nil {
		t.Fatal(err)
	}
}

func isolateHelperHooks(t *testing.T) {
	t.Helper()
	originalOpen, originalPipe, originalRandom := openHelperBinary, newPipeSFTP, randomHelperBytes
	t.Cleanup(func() { openHelperBinary, newPipeSFTP, randomHelperBytes = originalOpen, originalPipe, originalRandom })
	randomHelperBytes = func(data []byte) (int, error) {
		for index := range data {
			data[index] = byte(index)
		}
		return len(data), nil
	}
}

type fakeHelperSession struct {
	stdinErr, stdoutErr, startErr, combinedErr, waitErr, closeErr error
	combinedData                                                  []byte
	combinedBlock, closeSignal                                    chan struct{}
	stdin                                                         io.Reader
	closed                                                        bool
}

func (s *fakeHelperSession) SetStdin(input io.Reader) { s.stdin = input }
func (s *fakeHelperSession) CombinedOutput(string) ([]byte, error) {
	if s.combinedBlock != nil {
		<-s.combinedBlock
	}
	if s.stdin != nil {
		_, _ = io.Copy(io.Discard, s.stdin)
	}
	return s.combinedData, s.combinedErr
}
func (s *fakeHelperSession) StdinPipe() (io.WriteCloser, error) {
	if s.stdinErr != nil {
		return nil, s.stdinErr
	}
	return &testWriteCloser{}, nil
}
func (s *fakeHelperSession) StdoutPipe() (io.Reader, error) {
	if s.stdoutErr != nil {
		return nil, s.stdoutErr
	}
	return bytes.NewReader(nil), nil
}
func (s *fakeHelperSession) Start(string) error { return s.startErr }
func (s *fakeHelperSession) Wait() error        { return s.waitErr }
func (s *fakeHelperSession) Close() error {
	if !s.closed && s.closeSignal != nil {
		close(s.closeSignal)
	}
	if !s.closed && s.combinedBlock != nil {
		close(s.combinedBlock)
	}
	s.closed = true
	return s.closeErr
}

type testWriteCloser struct{ bytes.Buffer }

func (*testWriteCloser) Close() error { return nil }

type testReadCloser struct {
	io.Reader
	closeErr error
}

func (r *testReadCloser) Close() error { return r.closeErr }
