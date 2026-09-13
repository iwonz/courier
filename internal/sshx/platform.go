package sshx

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	posixPlatformProbe   = "uname -s; uname -m"
	windowsPlatformProbe = `powershell -NoProfile -NonInteractive -Command "[System.Environment]::OSVersion.Platform; $env:PROCESSOR_ARCHITECTURE"`
	posixArchiverProbe   = "command -v tar 2>/dev/null; command -v bsdtar 2>/dev/null; command -v gtar 2>/dev/null"
	windowsArchiverProbe = `powershell -NoProfile -NonInteractive -Command "Get-Command tar.exe,bsdtar.exe,gtar.exe -ErrorAction SilentlyContinue | ForEach-Object Source"`
)

// Runner executes constant remote capability probes.
type Runner interface {
	Run(context.Context, string) ([]byte, error)
}

// Platform is a normalized remote operating system and architecture.
type Platform struct {
	OS   string
	Arch string
}

// Archiver is either a discovered remote executable or the built-in fallback.
type Archiver struct {
	Name    string
	BuiltIn bool
}

// DetectPlatform probes POSIX first, then Windows, and discovers archivers.
func DetectPlatform(ctx context.Context, runner Runner) (Platform, Archiver, error) {
	if runner == nil {
		return Platform{}, Archiver{}, errors.New("remote probe runner is required")
	}
	output, posixErr := runner.Run(ctx, posixPlatformProbe)
	platform, err := parsePlatform(string(output))
	probeCommand := posixArchiverProbe
	if posixErr != nil || err != nil {
		output, windowsErr := runner.Run(ctx, windowsPlatformProbe)
		if windowsErr != nil {
			return Platform{}, Archiver{}, errors.Join(posixErr, err, windowsErr)
		}
		platform, err = parsePlatform("windows\n" + lastField(string(output)))
		if err != nil {
			return Platform{}, Archiver{}, err
		}
		probeCommand = windowsArchiverProbe
	}
	availableOutput, _ := runner.Run(ctx, probeCommand)
	return platform, SelectArchiver(platform, strings.Fields(string(availableOutput))), nil
}

func parsePlatform(output string) (Platform, error) {
	fields := strings.Fields(output)
	if len(fields) < 2 {
		return Platform{}, fmt.Errorf("unrecognized platform probe output %q", strings.TrimSpace(output))
	}
	osName := strings.ToLower(fields[0])
	switch osName {
	case "darwin":
	case "linux":
	case "freebsd", "openbsd", "netbsd", "dragonfly":
	case "windows", "win32nt", "windows_nt":
		osName = "windows"
	default:
		return Platform{}, fmt.Errorf("unsupported remote OS %q", fields[0])
	}
	architecture := strings.ToLower(fields[1])
	switch architecture {
	case "x86_64", "x64", "amd64":
		architecture = "amd64"
	case "aarch64", "arm64":
		architecture = "arm64"
	default:
		return Platform{}, fmt.Errorf("unsupported remote architecture %q", fields[1])
	}
	return Platform{OS: osName, Arch: architecture}, nil
}

// SelectArchiver deterministically picks a platform-appropriate executable.
func SelectArchiver(platform Platform, available []string) Archiver {
	set := make(map[string]string, len(available))
	for _, item := range available {
		base := strings.ToLower(filepath.Base(strings.ReplaceAll(item, `\`, "/")))
		set[base] = item
	}
	preference := []string{"tar", "gtar", "bsdtar"}
	if platform.OS == "darwin" || strings.Contains(platform.OS, "bsd") || platform.OS == "dragonfly" {
		preference = []string{"bsdtar", "tar", "gtar"}
	}
	if platform.OS == "windows" {
		preference = []string{"tar.exe", "bsdtar.exe", "gtar.exe"}
	}
	for _, candidate := range preference {
		if original := set[candidate]; original != "" {
			return Archiver{Name: original}
		}
	}
	return Archiver{Name: "builtin", BuiltIn: true}
}

func lastField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[len(fields)-1]
}
