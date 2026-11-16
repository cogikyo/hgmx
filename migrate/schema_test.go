package migrate_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
)

type owner struct {
	ID   string
	Name string `gorm:"not null;unique"`
}

type item struct {
	ID      string
	OwnerID string `gorm:"not null;index"`
	Owner   *owner `gorm:"constraint:OnDelete:CASCADE"`
	Kind    string `gorm:"type:item_kind;not null;check:kind <> 'a'"`
	Qty     int    `gorm:"type:integer;not null;default:1"`
}

type itemNote struct {
	ID   string
	Note string
}

func (itemNote) TableName() string { return "items" }

func kinds(values ...string) []migrate.Enum {
	return []migrate.Enum{pg.Enum[string]{Name: "item_kind", Values: values}}
}

var ddl = regexp.MustCompile(`(?i)^\s*(CREATE|ALTER|DROP|COMMENT)\b`)

var addValue = regexp.MustCompile(`^ALTER TYPE "[a-z_]+" ADD VALUE IF NOT EXISTS '[^']*'$`)

type recorder struct {
	logger.Interface
	ddl []string
}

func (r *recorder) LogMode(logger.LogLevel) logger.Interface { return r }

func (r *recorder) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	if sql, _ := fc(); ddl.MatchString(sql) && !addValue.MatchString(sql) {
		r.ddl = append(r.ddl, sql)
	}
}

func up(t *testing.T, db *gorm.DB, schemas ...migrate.Schema) error {
	t.Helper()
	return migrate.Run(t.Context(), db, schemas, []string{"up"}, io.Discard)
}

func files(m map[string]string) fstest.MapFS {
	fsys := fstest.MapFS{}
	for name, sql := range m {
		fsys[name] = &fstest.MapFile{Data: []byte(sql)}
	}
	return fsys
}

func ledger(t *testing.T, db *gorm.DB) string {
	t.Helper()
	return scalar[string](t, db, "SELECT coalesce(string_agg(schema || ' ' || file || ' ' || sha256, ', ' ORDER BY schema, file), '') FROM hgmx_init")
}

func digest(sql string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(sql))) }

func TestFailedUp(t *testing.T) {
	db := open(t)
	if err := up(t, db, migrate.Schema{Name: "site", Enums: kinds("a"), Models: []any{&itemNote{}}, Init: files(map[string]string{"init/post/a.sql": "SELECT 1"})}); err != nil {
		t.Fatal(err)
	}
	before := ledger(t, db)

	err := up(t, db, migrate.Schema{
		Name:   "site",
		Enums:  kinds("a", "c"),
		Models: []any{&owner{}, &item{}},
		Init:   files(map[string]string{"init/post/a.sql": "SELECT 2", "init/post/b.sql": "SELECT 1/0"}),
	})
	if err == nil || !strings.Contains(err.Error(), "post/b.sql") {
		t.Fatalf("up: err %v, want post/b.sql's failure", err)
	}
	if !scalar[bool](t, db, "SELECT to_regclass('owners') IS NULL") {
		t.Error("failed up left the owners table")
	}
	if scalar[bool](t, db, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'items' AND column_name = 'owner_id')") {
		t.Error("failed up left the items.owner_id column")
	}
	if got := ledger(t, db); got != before {
		t.Errorf("failed up changed the ledger to %q, want %q", got, before)
	}
	if !scalar[bool](t, db, "SELECT 'c' = ANY (enum_range(NULL::item_kind)::text[])") {
		t.Error("failed up lost the enum label committed before it")
	}
}

func TestLedgerSnapshot(t *testing.T) {
	db := open(t)
	pre, post := "CREATE TABLE IF NOT EXISTS notes (id int)", "SELECT 1"
	init := files(map[string]string{"init/pre/notes.sql": pre, "init/post/check.sql": post})
	if err := up(t, db, migrate.Schema{Name: "auth"}, migrate.Schema{Name: "site", Init: init}); err != nil {
		t.Fatal(err)
	}
	if got, want := ledger(t, db), "site post/check.sql "+digest(post)+", site pre/notes.sql "+digest(pre); got != want {
		t.Fatalf("ledger %q, want %q", got, want)
	}
	delete(init, "init/pre/notes.sql")
	if err := up(t, db, migrate.Schema{Name: "site", Init: init}); err != nil {
		t.Fatal(err)
	}
	if got, want := ledger(t, db), "site post/check.sql "+digest(post); got != want {
		t.Errorf("ledger %q after removing pre/notes.sql, want %q", got, want)
	}
}

func TestCommitRejected(t *testing.T) {
	db := open(t)
	err := up(t, db, migrate.Schema{Name: "site", Init: files(map[string]string{"init/post/commit.sql": "SELECT 1; COMMIT"})})
	if err == nil || !strings.Contains(err.Error(), "post/commit.sql ended the migration transaction") {
		t.Fatalf("up: err %v, want post/commit.sql rejected", err)
	}
}

func TestRerun(t *testing.T) {
	db := open(t)
	site := migrate.Schema{Name: "site", Enums: kinds("a", "b"), Models: []any{&owner{}, &item{}}}
	r := &recorder{Interface: logger.Discard}
	db = db.Session(&gorm.Session{Logger: r})
	if err := up(t, db, site); err != nil {
		t.Fatal(err)
	}
	if len(r.ddl) == 0 {
		t.Fatal("first up issued no DDL")
	}
	r.ddl = nil
	if err := up(t, db, site); err != nil {
		t.Fatal(err)
	}
	if len(r.ddl) > 0 {
		t.Errorf("rerun issued %q", r.ddl)
	}
}
