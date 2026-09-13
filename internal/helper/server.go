// Package helper implements Courier's temporary, consent-gated remote helper.
package helper

import (
	"errors"
	"io"

	"github.com/pkg/sftp"
)

type sftpServer interface {
	Serve() error
	Close() error
}

var newSFTPServer = func(stream io.ReadWriteCloser) (sftpServer, error) {
	return sftp.NewServer(stream, sftp.WithAllocator())
}

type splitStream struct {
	io.Reader
	io.Writer
}

func (splitStream) Close() error { return nil }

// ServeSFTP serves the established SFTP protocol over process standard streams.
func ServeSFTP(input io.Reader, output io.Writer) error {
	if input == nil || output == nil {
		return errors.New("helper standard streams are required")
	}
	server, err := newSFTPServer(splitStream{Reader: input, Writer: output})
	if err != nil {
		return err
	}
	serveErr := server.Serve()
	if errors.Is(serveErr, io.EOF) {
		serveErr = nil
	}
	return errors.Join(serveErr, server.Close())
}
