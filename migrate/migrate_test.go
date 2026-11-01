package migrate_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"slices"
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
	err := migrate.Command(t.Context(), db, []migrate.Track{{Steps: steps}}, args, &out)
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
	if !strings.Contains(out, "up 1 notes") {
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
		"status": "track site\nno migrations\n",
		"down":   "nothing to roll back\n",
	} {
		out, err := command(t, db, nil, arg)
		if err != nil || out != want {
			t.Errorf("%s: %q, %v", arg, out, err)
		}
	}
}

func TestTracks(t *testing.T) {
	db := open(t)
	steps := []migrate.Step{{
		Version: 1, Name: "notes",
		Up: func(ctx context.Context, db *gorm.DB) error {
			return db.Exec("CREATE TABLE notes (id int PRIMARY KEY)").Error
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return db.Exec("DROP TABLE notes").Error
		},
	}}
	if _, err := command(t, db, steps, "up"); err != nil {
		t.Fatal(err)
	}
	history := scalar[string](t, db, "SELECT json_agg(v ORDER BY id)::text FROM goose_db_version v WHERE version_id <= 1")
	var order []string
	var lockPIDs []int
	const lockQuery = "SELECT CASE WHEN count(*) = 1 THEN min(pid) ELSE 0 END FROM pg_locks WHERE locktype = 'advisory' AND database = (SELECT oid FROM pg_database WHERE datname = current_database()) AND classid = 0 AND objid = 4097083626 AND granted"
	tracks := []migrate.Track{
		{Name: "auth", Steps: []migrate.Step{{
			Version: 1, Name: "users",
			Up: func(ctx context.Context, db *gorm.DB) error {
				order = append(order, "auth")
				lockPIDs = append(lockPIDs, scalar[int](t, db, lockQuery))
				return db.Exec("CREATE TABLE auth_users (id int PRIMARY KEY)").Error
			},
			Down: func(ctx context.Context, db *gorm.DB) error {
				return db.Exec("DROP TABLE auth_users").Error
			},
		}}},
		{Steps: append(steps, migrate.Step{
			Version: 2, Name: "owners",
			Up: func(ctx context.Context, db *gorm.DB) error {
				order = append(order, "site")
				lockPIDs = append(lockPIDs, scalar[int](t, db, lockQuery))
				return db.Exec("ALTER TABLE notes ADD COLUMN user_id int REFERENCES auth_users(id)").Error
			},
			Down: func(ctx context.Context, db *gorm.DB) error {
				return db.Exec("ALTER TABLE notes DROP COLUMN user_id").Error
			},
		})},
	}
	run := func(args ...string) string {
		t.Helper()
		var out bytes.Buffer
		if err := migrate.Command(t.Context(), db, tracks, args, &out); err != nil {
			t.Fatal(err)
		}
		return out.String()
	}
	run("up")
	if !slices.Equal(order, []string{"auth", "site"}) {
		t.Fatalf("applied %v, want auth then site without replaying site step 1", order)
	}
	if len(lockPIDs) != 2 || lockPIDs[0] == 0 || lockPIDs[0] != lockPIDs[1] {
		t.Errorf("track lock sessions %v, want the same session holding one lock", lockPIDs)
	}
	if got := scalar[string](t, db, "SELECT json_agg(v ORDER BY id)::text FROM goose_db_version v WHERE version_id <= 1"); got != history {
		t.Errorf("existing site history changed: got %s, want %s", got, history)
	}
	if v := scalar[int](t, db, "SELECT max(version_id) FROM goose_db_version"); v != 2 {
		t.Errorf("site version %d, want 2", v)
	}
	if v := scalar[int](t, db, "SELECT max(version_id) FROM hgmx_auth_migrations"); v != 1 {
		t.Errorf("auth version %d, want 1", v)
	}
	if out := run("status"); !strings.Contains(out, "track auth") || !strings.Contains(out, "track site") || !strings.Contains(out, "owners") || !strings.Contains(out, "users") {
		t.Errorf("status omitted a track: %q", out)
	}
	run("up")
	if len(order) != 2 {
		t.Errorf("second up replayed steps: %v", order)
	}
	run("down")
	if v := scalar[int](t, db, "SELECT max(version_id) FROM goose_db_version"); v != 1 {
		t.Errorf("site down left version %d, want 1", v)
	}
	if v := scalar[int](t, db, "SELECT max(version_id) FROM hgmx_auth_migrations"); v != 1 {
		t.Errorf("site down changed auth version to %d", v)
	}
	if scalar[bool](t, db, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'notes' AND column_name = 'user_id')") {
		t.Error("site down left the user_id column")
	}
	run("down", "auth")
	if !scalar[bool](t, db, "SELECT to_regclass('auth_users') IS NULL") {
		t.Error("auth down left the auth_users table")
	}
	if v := scalar[int](t, db, "SELECT max(version_id) FROM hgmx_auth_migrations"); v != 0 {
		t.Errorf("auth down left version %d, want 0", v)
	}
	if got := scalar[string](t, db, "SELECT json_agg(v ORDER BY id)::text FROM goose_db_version v"); got != history {
		t.Errorf("auth down changed site history: got %s, want %s", got, history)
	}
}

func TestTrackValidation(t *testing.T) {
	t.Setenv("PGDATABASE", "")
	step := migrate.Step{Version: 1, Up: func(context.Context, *gorm.DB) error { return nil }}
	for _, tt := range []struct {
		name   string
		tracks []migrate.Track
		args   []string
		want   string
	}{
		{"duplicate", []migrate.Track{{Name: "auth"}, {Name: "auth"}}, []string{"up"}, "duplicate track"},
		{"duplicate site", []migrate.Track{{}, {}}, []string{"up"}, "duplicate track"},
		{"unsafe name", []migrate.Track{{Name: "auth.users"}}, []string{"up"}, "invalid track name"},
		{"long name", []migrate.Track{{Name: strings.Repeat("a", 48)}}, []string{"up"}, "invalid track name"},
		{"version", []migrate.Track{{Name: "auth", Steps: []migrate.Step{step}}, {Steps: []migrate.Step{{Version: 2, Up: step.Up}}}}, []string{"up"}, "has version 2"},
		{"nil up", []migrate.Track{{Name: "auth", Steps: []migrate.Step{{Version: 1}}}}, []string{"up"}, "has no Up"},
		{"unknown", nil, []string{"down", "auth"}, "unknown track"},
		{"up target", nil, []string{"up", "auth"}, "usage:"},
		{"named guard", []migrate.Track{{}, {Name: "auth", Steps: []migrate.Step{step}}}, []string{"status"}, "PGDATABASE is unset"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := migrate.Command(t.Context(), nil, tt.tracks, tt.args, &out)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v, want %q", err, tt.want)
			}
		})
	}
}
