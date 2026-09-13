package helper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/update"
	"github.com/pkg/sftp"
)

// Artifact is one locally verified helper executable.
type Artifact struct {
	Path    string
	SHA256  string
	Cleanup func() error
}

// Confirm asks for explicit consent before helper acquisition or deployment.
type Confirm func(context.Context, string) (bool, error)

// Acquire obtains an exact-platform helper without touching the remote host.
type Acquire func(context.Context, sshx.Platform) (*Artifact, error)

// Deploy stages and starts a helper over an authenticated SSH connection.
type Deploy func(context.Context, *sshx.Connection, sshx.Platform, string, string) (*sftp.Client, io.Closer, error)

// Detect identifies the authenticated remote platform without changing it.
type Detect func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error)

// Manager applies the capability, consent, acquisition, and cleanup policy.
type Manager struct {
	Confirm Confirm
	Acquire Acquire
	Deploy  Deploy
	Detect  Detect
	mu      sync.Mutex
}

var detectRemotePlatform = sshx.DetectPlatform

// Fallback satisfies sshx.Factory.SFTPFallback.
func (m *Manager) Fallback(ctx context.Context, connection *sshx.Connection, capability *sshx.CapabilityError) (*sftp.Client, io.Closer, error) {
	if capability == nil || capability.Capability != "sftp" {
		return nil, nil, errors.New("remote helper requires an SFTP capability error")
	}
	if m == nil || m.Confirm == nil {
		return nil, nil, fmt.Errorf("%w; interactive consent is unavailable", capability)
	}
	detect := m.Detect
	if detect == nil {
		detect = detectRemotePlatform
	}
	platform, _, err := detect(ctx, connection)
	if err != nil {
		return nil, nil, fmt.Errorf("detect helper platform: %w", err)
	}
	m.mu.Lock()
	accepted, confirmErr := m.Confirm(ctx, fmt.Sprintf("Remote %s lacks SFTP. Deploy a temporary verified Courier helper for %s/%s?", connection.Target.Alias, platform.OS, platform.Arch))
	m.mu.Unlock()
	if confirmErr != nil {
		return nil, nil, fmt.Errorf("%w; helper consent failed: %v", capability, confirmErr)
	}
	if !accepted {
		return nil, nil, fmt.Errorf("%w; helper deployment declined", capability)
	}
	if m.Acquire == nil || m.Deploy == nil {
		return nil, nil, errors.New("remote helper dependencies are incomplete")
	}
	artifact, err := m.Acquire(ctx, platform)
	if err != nil {
		return nil, nil, fmt.Errorf("acquire remote helper: %w", err)
	}
	if artifact == nil || artifact.Path == "" || artifact.Cleanup == nil {
		return nil, nil, errors.New("remote helper acquisition returned an incomplete artifact")
	}
	client, helperRuntime, deployErr := m.Deploy(ctx, connection, platform, artifact.Path, artifact.SHA256)
	if deployErr == nil && (client == nil || helperRuntime == nil) {
		deployErr = errors.New("remote helper deployment returned an incomplete runtime")
	}
	cleanupErr := artifact.Cleanup()
	if deployErr != nil || cleanupErr != nil {
		if client != nil {
			_ = client.Close()
		}
		if helperRuntime != nil {
			_ = helperRuntime.Close()
		}
		return nil, nil, errors.Join(deployErr, cleanupErr)
	}
	return client, helperRuntime, nil
}

// Source selects the running executable for matching platforms and immutable
// GitHub Release artifacts for cross-platform helpers.
type Source struct {
	Repository string
	Version    string
	Token      string
	Client     update.HTTPClient
	GOOS       string
	GOARCH     string
	Executable func() (string, error)
	EvalLinks  func(string) (string, error)
}

var openArtifactFile = func(name string) (io.ReadCloser, error) { return os.Open(name) }

// Acquire returns a verified helper for platform.
func (s Source) Acquire(ctx context.Context, platform sshx.Platform) (*Artifact, error) {
	if !supportedHelperPlatform(platform) {
		return nil, fmt.Errorf("no Courier helper release for %s/%s", platform.OS, platform.Arch)
	}
	goos, goarch := s.GOOS, s.GOARCH
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	if platform.OS == goos && platform.Arch == goarch {
		executable := s.Executable
		if executable == nil {
			executable = os.Executable
		}
		name, err := executable()
		if err != nil {
			return nil, err
		}
		evaluate := s.EvalLinks
		if evaluate == nil {
			evaluate = filepath.EvalSymlinks
		}
		name, err = evaluate(name)
		if err != nil {
			return nil, err
		}
		digest, err := artifactDigest(name)
		if err != nil {
			return nil, err
		}
		return &Artifact{Path: name, SHA256: digest, Cleanup: func() error { return nil }}, nil
	}
	updater := update.Updater{Repository: s.Repository, Version: s.Version, Token: s.Token, Client: s.Client}
	binary, err := updater.AcquireVersion(ctx, s.Version, platform.OS, platform.Arch)
	if err != nil {
		return nil, err
	}
	return &Artifact{Path: binary.Path, SHA256: binary.SHA256, Cleanup: binary.Cleanup}, nil
}

func supportedHelperPlatform(platform sshx.Platform) bool {
	if platform.Arch != "amd64" && platform.Arch != "arm64" {
		return false
	}
	switch strings.ToLower(platform.OS) {
	case "darwin", "linux", "windows", "freebsd", "openbsd", "netbsd":
		return true
	case "dragonfly":
		return platform.Arch == "amd64"
	default:
		return false
	}
}

func artifactDigest(name string) (string, error) {
	file, err := openArtifactFile(name)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	_, copyErr := io.Copy(hash, file)
	closeErr := file.Close()
	if copyErr != nil || closeErr != nil {
		return "", errors.Join(copyErr, closeErr)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
