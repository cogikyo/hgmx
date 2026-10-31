// Package pgtest supplies isolated databases on a shared local Postgres cluster.
package pgtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/cogikyo/hgmx/pg/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type server struct {
	dir   string
	admin *sql.DB
}

var cluster = sync.OnceValues(start)

// Open creates a test database and non-superuser role, cleaned up with t.
// Cluster startup errors skip the test unless HGMX_PGTEST=require.
// The per-user cluster lives in $TMPDIR/hgmx-pgtest-<uid> and stays running.
func Open(t testing.TB) *gorm.DB {
	t.Helper()
	s, err := cluster()
	if err != nil {
		if os.Getenv("HGMX_PGTEST") == "require" {
			t.Fatalf("pgtest: %v", err)
		}
		t.Skipf("pgtest: %v", err)
	}
	name := fmt.Sprintf("t%016x", rand.Uint64())
	admin := func(stmt string) error {
		_, err := s.admin.ExecContext(context.Background(), stmt)
		return err
	}
	if err := admin("CREATE ROLE " + name + " LOGIN"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin("DROP ROLE " + name); err != nil {
			t.Error(err)
		}
	})
	if err := admin("CREATE DATABASE " + name + " OWNER " + name); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin("DROP DATABASE " + name + " WITH (FORCE)"); err != nil {
			t.Error(err)
		}
	})
	log := slog.New(slog.NewTextHandler(t.Output(), nil))
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn(s.dir, name)}), config.Gorm(log))
	if err != nil {
		t.Fatal(err)
	}
	conn, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return db
}

func start() (*server, error) {
	for _, bin := range []string{"initdb", "pg_ctl"} {
		if _, err := exec.LookPath(bin); err != nil {
			return nil, err
		}
	}
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("hgmx-pgtest-%d", os.Getuid()))
	if socket := filepath.Join(dir, ".s.PGSQL.5432"); len(socket) > 107 {
		return nil, fmt.Errorf("socket path %s is %d bytes, over the 107-byte limit; set TMPDIR to a shorter directory", socket, len(socket))
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return nil, err
	}

	data := filepath.Join(dir, "data")
	if _, err := os.Stat(filepath.Join(data, "PG_VERSION")); errors.Is(err, fs.ErrNotExist) {
		if err := initdb(dir, data); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	if exec.Command("pg_ctl", "status", "-D", data).Run() != nil {
		logfile := filepath.Join(dir, "log")
		cmd := exec.Command("pg_ctl", "start", "-w", "-s", "-D", data, "-l", logfile)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if out, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("pg_ctl start: %v: %s (server log: %s)", err, out, logfile)
		}
	}
	admin, err := sql.Open("pgx", dsn(dir, "postgres"))
	if err != nil {
		return nil, err
	}
	return &server{dir: dir, admin: admin}, nil
}

func initdb(dir, data string) error {
	tmp := data + ".tmp"
	if err := os.RemoveAll(tmp); err != nil {
		return err
	}
	out, err := exec.Command("initdb", "-D", tmp, "-U", "postgres", "-A", "trust", "-E", "UTF8", "--no-locale", "-N").CombinedOutput()
	if err != nil {
		return fmt.Errorf("initdb: %v: %s", err, out)
	}
	conf, err := os.OpenFile(filepath.Join(tmp, "postgresql.conf"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(conf, "listen_addresses = ''\nunix_socket_directories = '%s'\nfsync = off\n", strings.ReplaceAll(dir, "'", "''"))
	if err := errors.Join(err, conf.Close()); err != nil {
		return err
	}
	return os.Rename(tmp, data)
}

func dsn(dir, name string) string {
	quote := strings.NewReplacer(`\`, `\\`, `'`, `\'`)
	return fmt.Sprintf("host='%s' port=5432 user=%s dbname=%s sslmode=disable", quote.Replace(dir), name, name)
}
