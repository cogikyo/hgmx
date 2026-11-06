package main

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	bufPin         = "github.com/bufbuild/buf/cmd/buf@v1.75.0"
	protocGenGo    = "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12"
	protocGenGoRT  = "v1.36.12"
	connectGo      = "connectrpc.com/connect/v2/cmd/protoc-gen-connect-go@v2.0.2"
	connectGoRT    = "v2.0.2"
	protobufGoPath = "google.golang.org/protobuf"
	connectGoPath  = "connectrpc.com/connect/v2"
)

type download struct {
	name, version, url, sum string
}

var protocs = map[string]download{
	"linux/amd64":  {"protoc", "36.2", "https://github.com/protocolbuffers/protobuf/releases/download/v36.2/protoc-36.2-linux-x86_64.zip", "121f6c7afe1d4d0e3ea6aab9432038599250134cbf4474cb1167d2c7decd4278"},
	"linux/arm64":  {"protoc", "36.2", "https://github.com/protocolbuffers/protobuf/releases/download/v36.2/protoc-36.2-linux-aarch_64.zip", "8b8f18bd2b30346efbc698dd5a73dd7c805f3ef8380f6dfc95c768f3f1852f6a"},
	"darwin/arm64": {"protoc", "36.2", "https://github.com/protocolbuffers/protobuf/releases/download/v36.2/protoc-36.2-osx-aarch_64.zip", "9cd98a532c5c5e0c4161314de0225de27e4c8a323917b6ea7b1b714d3ae23466"},
	"darwin/amd64": {"protoc", "36.2", "https://github.com/protocolbuffers/protobuf/releases/download/v36.2/protoc-36.2-osx-x86_64.zip", "228cc7add4616cc14ca5e80dee83209d44449a7aee95a914ae748fa374efb078"},
}

var connectKotlin = download{"connect-kotlin", "0.10.0", "https://repo1.maven.org/maven2/com/connectrpc/protoc-gen-connect-kotlin/0.10.0/protoc-gen-connect-kotlin-0.10.0.jar", "30d166ba87b6c9a4dfa501a2241f23b0c6523960828a798bce45294988ccc65b"}

type swiftTool struct {
	name, version, commit, url, product string
}

var swiftTools = []swiftTool{
	{"swift-protobuf", "1.39.0", "3d537cb81aef272a03cb3be815439e9bf0eb3bd3", "https://github.com/apple/swift-protobuf.git", "protoc-gen-swift"},
	{"connect-swift", "1.2.3", "ed56816fe90ab872566f09ed73e28ae078057e3d", "https://github.com/connectrpc/connect-swift.git", "protoc-gen-connect-swift"},
}

//go:embed swifttools.sh
var swiftScript []byte

type tools struct {
	log *slog.Logger
}

func toolDir(name, version string) (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, "hgmx", "tools", name, version), nil
}

func toolEnv() []string {
	return append(os.Environ(), "GOWORK=off", "GOFLAGS=", "CGO_ENABLED=0", "GOTOOLCHAIN=local", "GIT_TERMINAL_PROMPT=0")
}

func bufCmd(ctx context.Context, dir string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", bufPin}, args...)...)
	cmd.Dir = dir
	cmd.Env = toolEnv()
	cmd.WaitDelay = time.Second
	return cmd
}

func bufPassCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cmd := bufCmd(ctx, ".", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	if err != nil {
		log.Error("Failed to run buf", slog.String("error", err.Error()))
		return 1
	}
	return 0
}

func (t tools) protoc(ctx context.Context) (string, error) {
	d, ok := protocs[runtime.GOOS+"/"+runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("hgmx pins no protoc for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	dir, err := t.fetch(ctx, d, func(data []byte, dir string) error {
		r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return err
		}
		return os.CopyFS(dir, r)
	})
	return filepath.Join(dir, "bin", "protoc"), err
}

func (t tools) connectKotlin(ctx context.Context) (string, error) {
	name := filepath.Base(connectKotlin.url)
	dir, err := t.fetch(ctx, connectKotlin, func(data []byte, dir string) error {
		return os.WriteFile(filepath.Join(dir, name), data, 0o644)
	})
	return filepath.Join(dir, name), err
}

func (t tools) fetch(ctx context.Context, d download, unpack func(data []byte, dir string) error) (string, error) {
	dir, err := toolDir(d.name, d.version)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err == nil {
		t.log.Info("Tool cached", slog.String("tool", d.name), slog.String("version", d.version))
		return dir, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download %s %s: %w", d.name, d.version, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s %s: %s", d.name, d.version, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("download %s %s: %w", d.name, d.version, err)
	}
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != d.sum {
		return "", fmt.Errorf("download %s %s: sha256 %x, want %s", d.name, d.version, sum, d.sum)
	}
	parent := filepath.Dir(dir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(parent, "."+d.version+".hgmx-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	if err := unpack(data, tmp); err != nil {
		return "", fmt.Errorf("unpack %s %s: %w", d.name, d.version, err)
	}
	if err := os.Rename(tmp, dir); err != nil {
		if _, serr := os.Stat(dir); serr != nil {
			return "", err
		}
	}
	t.log.Info("Tool installed", slog.String("tool", d.name), slog.String("version", d.version), slog.String("dir", dir))
	return dir, nil
}

type mac struct {
	host string
}

func macHost() mac {
	if runtime.GOOS == "darwin" {
		return mac{}
	}
	host := os.Getenv("HGMX_MAC")
	if host == "" {
		host = "neumann"
	}
	return mac{host: host}
}

func (m mac) String() string {
	if m.host == "" {
		return "this Mac"
	}
	return m.host
}

func (m mac) command(ctx context.Context, args ...string) *exec.Cmd {
	if m.host == "" {
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		cmd.Dir = os.Getenv("HOME")
		return cmd
	}
	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
	}
	return exec.CommandContext(ctx, "ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=10", m.host, strings.Join(quoted, " "))
}

func (m mac) plugin(t swiftTool) []string {
	rel := filepath.Join("Library", "Caches", "hgmx", "tools", t.name, t.version, t.product)
	if m.host == "" {
		return []string{filepath.Join(os.Getenv("HOME"), rel)}
	}
	return []string{"ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=10", m.host, rel}
}

func (t tools) swift(ctx context.Context, m mac) error {
	for _, st := range swiftTools {
		cmd := m.command(ctx, "sh", "-s", "--", st.name, st.version, st.commit, st.url, st.product)
		cmd.Stdin = bytes.NewReader(swiftScript)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("%s %s on %s: %w\n%s", st.product, st.version, m, err, tail(stderr.String(), 20))
		}
		if strings.TrimSpace(string(out)) == "cached" {
			t.log.Info("Tool cached", slog.String("tool", st.product), slog.String("version", st.version), slog.String("host", m.String()))
		} else {
			t.log.Info("Tool built", slog.String("tool", st.product), slog.String("version", st.version), slog.String("host", m.String()))
		}
	}
	return nil
}

func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
