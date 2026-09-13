package helper

import (
	"errors"
	"io"
	"net"
	"testing"

	"github.com/pkg/sftp"
)

func TestServeSFTP(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	done := make(chan error, 1)
	go func() { done <- ServeSFTP(serverSide, serverSide) }()
	client, err := sftp.NewClientPipe(clientSide, clientSide)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Lstat(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestServeSFTPFailures(t *testing.T) {
	if err := ServeSFTP(nil, io.Discard); err == nil {
		t.Fatal("expected missing stream error")
	}
	if err := ServeSFTP(&emptyReader{}, nil); err == nil {
		t.Fatal("expected missing output error")
	}
	original := newSFTPServer
	t.Cleanup(func() { newSFTPServer = original })
	newSFTPServer = func(io.ReadWriteCloser) (sftpServer, error) { return nil, errors.New("create") }
	if err := ServeSFTP(&memoryStream{}, io.Discard); err == nil {
		t.Fatal("expected create error")
	}
	serveCause := errors.New("serve")
	newSFTPServer = func(io.ReadWriteCloser) (sftpServer, error) {
		return fakeSFTPServer{serveErr: serveCause, closeErr: errors.New("close")}, nil
	}
	if err := ServeSFTP(&memoryStream{}, io.Discard); err == nil || !errors.Is(err, serveCause) {
		t.Fatal("expected joined server errors")
	}
	newSFTPServer = func(io.ReadWriteCloser) (sftpServer, error) {
		return fakeSFTPServer{serveErr: io.EOF}, nil
	}
	if err := ServeSFTP(&memoryStream{}, io.Discard); err != nil {
		t.Fatalf("EOF should be a clean shutdown: %v", err)
	}
}

type fakeSFTPServer struct{ serveErr, closeErr error }

func (f fakeSFTPServer) Serve() error { return f.serveErr }
func (f fakeSFTPServer) Close() error { return f.closeErr }

type emptyReader struct{}

func (*emptyReader) Read([]byte) (int, error) { return 0, io.EOF }

type memoryStream struct{ emptyReader }

func (*memoryStream) Write(data []byte) (int, error) { return len(data), nil }
func (*memoryStream) Close() error                   { return nil }
