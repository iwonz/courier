package sshx

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var helperDigestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

type helperSession interface {
	SetStdin(io.Reader)
	CombinedOutput(string) ([]byte, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
	Start(string) error
	Wait() error
	Close() error
}

type realHelperSession struct{ session *ssh.Session }

func (s realHelperSession) SetStdin(input io.Reader) { s.session.Stdin = input }
func (s realHelperSession) CombinedOutput(command string) ([]byte, error) {
	return s.session.CombinedOutput(command)
}
func (s realHelperSession) StdinPipe() (io.WriteCloser, error) { return s.session.StdinPipe() }
func (s realHelperSession) StdoutPipe() (io.Reader, error)     { return s.session.StdoutPipe() }
func (s realHelperSession) Start(command string) error         { return s.session.Start(command) }
func (s realHelperSession) Wait() error                        { return s.session.Wait() }
func (s realHelperSession) Close() error                       { return s.session.Close() }

var (
	openHelperBinary  = func(name string) (io.ReadCloser, error) { return os.Open(name) }
	newPipeSFTP       = sftp.NewClientPipe
	randomHelperBytes = rand.Read
)

func (c *Connection) newHelperSession() (helperSession, error) {
	if c != nil && c.helperSessionFactory != nil {
		return c.helperSessionFactory()
	}
	if c == nil || c.SSH == nil {
		return nil, errors.New("SSH connection is closed")
	}
	session, err := c.SSH.NewSession()
	if err != nil {
		return nil, err
	}
	return realHelperSession{session: session}, nil
}

// DeployHelper uploads a verified Courier binary and starts its embedded SFTP
// server. The returned closer removes the remote staging directory.
func DeployHelper(ctx context.Context, connection *Connection, platform Platform, binaryPath, digest string) (*sftp.Client, io.Closer, error) {
	if !helperDigestPattern.MatchString(digest) {
		return nil, nil, errors.New("helper SHA-256 digest is invalid")
	}
	remotePath, err := newRemoteHelperPath(platform)
	if err != nil {
		return nil, nil, err
	}
	upload, err := connection.newHelperSession()
	if err != nil {
		return nil, nil, err
	}
	defer upload.Close()
	binary, err := openHelperBinary(binaryPath)
	if err != nil {
		return nil, nil, err
	}
	upload.SetStdin(binary)
	_, uploadErr := runHelperCommand(ctx, upload, helperUploadCommand(platform, remotePath, digest))
	closeBinaryErr := binary.Close()
	if uploadErr != nil || closeBinaryErr != nil {
		runtime := &helperRuntime{connection: connection, platform: platform, remotePath: remotePath}
		return nil, nil, errors.Join(uploadErr, closeBinaryErr, runtime.cleanup())
	}
	return startHelper(ctx, connection, platform, remotePath)
}

func runHelperCommand(ctx context.Context, session helperSession, command string) ([]byte, error) {
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

func startHelper(ctx context.Context, connection *Connection, platform Platform, remotePath string) (*sftp.Client, io.Closer, error) {
	runtime := &helperRuntime{connection: connection, platform: platform, remotePath: remotePath}
	session, err := connection.newHelperSession()
	if err != nil {
		return nil, nil, errors.Join(err, runtime.cleanup())
	}
	runtime.session = session
	input, err := session.StdinPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, errors.Join(err, runtime.cleanup())
	}
	output, err := session.StdoutPipe()
	if err != nil {
		_ = input.Close()
		_ = session.Close()
		return nil, nil, errors.Join(err, runtime.cleanup())
	}
	if err := session.Start(helperStartCommand(platform, remotePath)); err != nil {
		_ = input.Close()
		_ = session.Close()
		return nil, nil, errors.Join(err, runtime.cleanup())
	}
	type clientOutcome struct {
		client *sftp.Client
		err    error
	}
	done := make(chan clientOutcome, 1)
	go func() {
		client, err := newPipeSFTP(output, input)
		done <- clientOutcome{client: client, err: err}
	}()
	select {
	case <-ctx.Done():
		_ = session.Close()
		return nil, nil, errors.Join(ctx.Err(), runtime.cleanup())
	case result := <-done:
		if result.err != nil {
			_ = session.Close()
			return nil, nil, errors.Join(result.err, runtime.cleanup())
		}
		return result.client, runtime, nil
	}
}

type helperRuntime struct {
	connection *Connection
	session    helperSession
	platform   Platform
	remotePath string
	once       sync.Once
	err        error
}

func (r *helperRuntime) Close() error {
	r.once.Do(func() {
		var waitErr error
		if r.session != nil {
			waitErr = r.session.Wait()
			_ = r.session.Close()
		}
		r.err = errors.Join(waitErr, r.cleanup())
	})
	return r.err
}

func (r *helperRuntime) cleanup() error {
	if r.connection == nil || r.remotePath == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := r.connection.Run(ctx, helperCleanupCommand(r.platform, r.remotePath))
	return err
}

func helperUploadCommand(platform Platform, remotePath, digest string) string {
	if platform.OS == "windows" {
		directory := strings.Split(strings.ReplaceAll(remotePath, `/`, `\`), `\`)[0]
		script := `$ErrorActionPreference='Stop';$d=Join-Path ([IO.Path]::GetTempPath()) ` + powerShellQuote(directory) + `;$p=Join-Path $d 'courier.exe';try{[IO.Directory]::CreateDirectory($d)|Out-Null;$i=[Console]::OpenStandardInput();$o=[IO.File]::Open($p,[IO.FileMode]::CreateNew,[IO.FileAccess]::Write,[IO.FileShare]::None);try{$i.CopyTo($o)}finally{$o.Dispose()};$h=(Get-FileHash -Algorithm SHA256 -LiteralPath $p).Hash.ToLowerInvariant();if($h -ne '` + digest + `'){throw 'helper checksum mismatch'}}catch{Remove-Item -LiteralPath $d -Recurse -Force -ErrorAction SilentlyContinue;[Console]::Error.WriteLine($_.Exception.Message);exit 1}`
		return encodedPowerShell(script)
	}
	directory := path.Dir(remotePath)
	return `umask 077; d=` + shellQuote(directory) + `; p=` + shellQuote(remotePath) + `; mkdir -m 700 "$d" || exit 70; if ! cat >"$p" || ! chmod 700 "$p"; then rm -rf "$d"; exit 71; fi; if command -v sha256sum >/dev/null 2>&1; then h=$(sha256sum "$p"); h=${h%% *}; elif command -v shasum >/dev/null 2>&1; then h=$(shasum -a 256 "$p"); h=${h%% *}; elif command -v sha256 >/dev/null 2>&1; then h=$(sha256 -q "$p"); else rm -rf "$d"; exit 72; fi; if [ "$h" != '` + digest + `' ]; then rm -rf "$d"; exit 73; fi`
}

func helperStartCommand(platform Platform, remotePath string) string {
	if platform.OS == "windows" {
		directory := strings.Split(strings.ReplaceAll(remotePath, `/`, `\`), `\`)[0]
		script := `$p=Join-Path (Join-Path ([IO.Path]::GetTempPath()) ` + powerShellQuote(directory) + `) 'courier.exe';& $p _helper-sftp`
		return encodedPowerShell(script)
	}
	return "exec " + shellQuote(remotePath) + " _helper-sftp"
}

func helperCleanupCommand(platform Platform, remotePath string) string {
	if platform.OS == "windows" {
		directory := strings.Split(strings.ReplaceAll(remotePath, `/`, `\`), `\`)[0]
		return encodedPowerShell(`$d=Join-Path ([IO.Path]::GetTempPath()) ` + powerShellQuote(directory) + `;Remove-Item -LiteralPath $d -Recurse -Force -ErrorAction Stop`)
	}
	return "rm -rf -- " + shellQuote(path.Dir(remotePath))
}

func newRemoteHelperPath(platform Platform) (string, error) {
	data := make([]byte, 16)
	if _, err := randomHelperBytes(data); err != nil {
		return "", err
	}
	name := "courier-helper-" + hex.EncodeToString(data)
	if platform.OS == "windows" {
		return name + `\courier.exe`, nil
	}
	return "/tmp/" + name + "/courier", nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func powerShellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func encodedPowerShell(script string) string {
	data := make([]byte, 0, len(script)*2)
	for _, character := range []rune(script) {
		if character <= 0xffff {
			data = append(data, byte(character), byte(character>>8))
			continue
		}
		character -= 0x10000
		high := rune(0xd800) + (character >> 10)
		low := rune(0xdc00) + (character & 0x3ff)
		data = append(data, byte(high), byte(high>>8), byte(low), byte(low>>8))
	}
	return "powershell -NoProfile -NonInteractive -EncodedCommand " + base64.StdEncoding.EncodeToString(data)
}
