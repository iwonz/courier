package installtest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPOSIXInstallerSuccessLatestAndCleanup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the POSIX installer runs on Unix hosts")
	}
	for _, version := range []string{"1.2.3", "latest"} {
		t.Run(version, func(t *testing.T) {
			fixture, scratch, installDir, asset := installerFixture(t, []byte("standalone-binary"), false)
			environment := installerEnvironment(fixture, scratch, installDir)
			if version == "latest" {
				releaseJSON := filepath.Join(fixture, "latest.json")
				if err := os.WriteFile(releaseJSON, []byte(`{"tag_name":"v1.2.3"}`), 0o600); err != nil {
					t.Fatal(err)
				}
				environment = append(environment, "COURIER_VERSION=latest", "COURIER_RELEASE_API_URL=file://"+releaseJSON)
			}
			command := exec.Command("sh", filepath.Join(projectRoot(t), "install.sh"))
			command.Env = environment
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("installer failed: %v\n%s", err, output)
			}
			installed, err := os.ReadFile(filepath.Join(installDir, "courier"))
			if err != nil || string(installed) != "standalone-binary" {
				t.Fatalf("installed=%q err=%v", installed, err)
			}
			info, err := os.Stat(filepath.Join(installDir, "courier"))
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != 0o755 {
				t.Fatalf("mode=%v", info.Mode())
			}
			assertEmpty(t, scratch)
			if !strings.Contains(string(output), asset) && !strings.Contains(string(output), "Installed courier 1.2.3") {
				t.Fatalf("unexpected output: %s", output)
			}
		})
	}
}

func TestPOSIXInstallerRejectsChecksumWithoutReplacing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the POSIX installer runs on Unix hosts")
	}
	fixture, scratch, installDir, _ := installerFixture(t, []byte("corrupt"), true)
	if err := os.MkdirAll(installDir, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(installDir, "courier")
	if err := os.WriteFile(target, []byte("existing"), 0o700); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("sh", filepath.Join(projectRoot(t), "install.sh"))
	command.Env = installerEnvironment(fixture, scratch, installDir)
	output, err := command.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "checksum mismatch") {
		t.Fatalf("expected checksum failure, err=%v output=%s", err, output)
	}
	installed, readErr := os.ReadFile(target)
	if readErr != nil || string(installed) != "existing" {
		t.Fatalf("existing target changed: %q err=%v", installed, readErr)
	}
	assertEmpty(t, scratch)
}

func installerFixture(t *testing.T, binary []byte, corruptChecksum bool) (fixture, scratch, installDir, asset string) {
	t.Helper()
	root := t.TempDir()
	fixture = filepath.Join(root, "releases")
	scratch = filepath.Join(root, "scratch")
	installDir = filepath.Join(root, "install")
	if err := os.MkdirAll(filepath.Join(fixture, "v1.2.3"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(scratch, 0o700); err != nil {
		t.Fatal(err)
	}
	osName := runtime.GOOS
	architecture := map[string]string{"amd64": "amd64", "arm64": "arm64"}[runtime.GOARCH]
	if (osName != "darwin" && osName != "linux") || architecture == "" {
		t.Skipf("unsupported installer test host %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	asset = fmt.Sprintf("courier_1.2.3_%s_%s", osName, architecture)
	assetPath := filepath.Join(fixture, "v1.2.3", asset)
	if err := os.WriteFile(assetPath, binary, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(binary)
	checksum := hex.EncodeToString(digest[:])
	if corruptChecksum {
		checksum = strings.Repeat("0", 64)
	}
	manifest := []byte(checksum + "  " + asset + "\n")
	if err := os.WriteFile(filepath.Join(fixture, "v1.2.3", "checksums.txt"), manifest, 0o600); err != nil {
		t.Fatal(err)
	}
	return fixture, scratch, installDir, asset
}

func installerEnvironment(fixture, scratch, installDir string) []string {
	return append(os.Environ(),
		"COURIER_VERSION=1.2.3",
		"COURIER_RELEASE_BASE_URL=file://"+fixture,
		"COURIER_INSTALL_DIR="+installDir,
		"TMPDIR="+scratch,
	)
}

func projectRoot(t *testing.T) string {
	t.Helper()
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate installer test")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(current), "..", ".."))
}

func assertEmpty(t *testing.T, directory string) {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 0 {
		t.Fatalf("temporary directory not cleaned: entries=%v err=%v", entries, err)
	}
}
