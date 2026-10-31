package migrate_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/pg/pgtest"
	"gorm.io/gorm"
)

type note struct {
	ID   int
	Body string
}

func open(t *testing.T) *gorm.DB {
	t.Setenv("PGDATABASE", "pgtest")
	return pgtest.Open(t)
}

func command(t *testing.T, db *gorm.DB, steps []migrate.Step, args ...string) (string, error) {
	var out bytes.Buffer
	err := migrate.Command(t.Context(), db, steps, args, &out)
	return out.String(), err
}

func scalar[T any](t *testing.T, db *gorm.DB, query string) T {
	var v T
	if err := db.Raw(query).Scan(&v).Error; err != nil {
		t.Fatal(err)
	}
	return v
}

func TestTx(t *testing.T) {
	db := open(t)
	steps := []migrate.Step{
		{
			Version: 1, Name: "notes",
			Up: func(ctx context.Context, db *gorm.DB) error {
				if err := db.Exec("CREATE TABLE notes (id int PRIMARY KEY, body text)").Error; err != nil {
					return err
				}
				var timeout string
				if err := db.Raw("SHOW lock_timeout").Scan(&timeout).Error; err != nil {
					return err
				}
				if timeout != "5s" {
					return errors.New("lock_timeout is " + timeout + ", want 5s inside the step")
				}
				return db.Create(&note{1, "kept"}).Error
			},
			Down: func(ctx context.Context, db *gorm.DB) error {
				return db.Exec("DROP TABLE notes").Error
			},
		},
		{
			Version: 2, Name: "broken",
			Up: func(ctx context.Context, db *gorm.DB) error {
				if err := db.Create(&note{2, "rolled back"}).Error; err != nil {
					return err
				}
				return errors.New("boom")
			},
		},
	}

	out, err := command(t, db, steps, "up")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("up: err %v, want boom", err)
	}
	if !strings.HasPrefix(out, "up 1 notes") {
		t.Errorf("up printed %q", out)
	}
	if n := scalar[int](t, db, "SELECT count(*) FROM notes"); n != 1 {
		t.Errorf("notes has %d rows, want step 1's only", n)
	}
	if v := scalar[int](t, db, "SELECT max(version_id) FROM goose_db_version"); v != 1 {
		t.Errorf("history at version %d, want 1", v)
	}

	if _, err := command(t, db, steps, "down"); err != nil {
		t.Fatal(err)
	}
	if !scalar[bool](t, db, "SELECT to_regclass('notes') IS NULL") {
		t.Error("down left the notes table")
	}
}

func TestNilDown(t *testing.T) {
	db := open(t)
	steps := []migrate.Step{{
		Version: 1, Name: "notes",
		Up: func(ctx context.Context, db *gorm.DB) error {
			return db.Exec("CREATE TABLE notes (id int)").Error
		},
	}}
	if _, err := command(t, db, steps, "up"); err != nil {
		t.Fatal(err)
	}
	_, err := command(t, db, steps, "down")
	if err == nil || !strings.Contains(err.Error(), "has no Down") {
		t.Fatalf("down: err %v, want has no Down", err)
	}
	if !scalar[bool](t, db, "SELECT EXISTS (SELECT 1 FROM goose_db_version WHERE version_id = 1)") {
		t.Error("history row for step 1 is gone")
	}
}

func TestEmpty(t *testing.T) {
	t.Setenv("PGDATABASE", "")
	db, err := pg.Open(slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	for arg, want := range map[string]string{
		"up":     "no migrations\n",
		"status": "no migrations\n",
		"down":   "nothing to roll back\n",
	} {
		out, err := command(t, db, nil, arg)
		if err != nil || out != want {
			t.Errorf("%s: %q, %v", arg, out, err)
		}
	}
}
