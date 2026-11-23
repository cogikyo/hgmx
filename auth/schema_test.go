package auth

import (
	"context"
	"io"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg/pgtest"
)

var ddl = regexp.MustCompile(`(?i)^\s*(CREATE|ALTER|DROP|COMMENT)\b`)

type recorder struct {
	logger.Interface
	ddl []string
}

func (r *recorder) LogMode(logger.LogLevel) logger.Interface { return r }

func (r *recorder) Trace(_ context.Context, _ time.Time, fc func() (string, int64), _ error) {
	if sql, _ := fc(); ddl.MatchString(sql) {
		r.ddl = append(r.ddl, sql)
	}
}

const catalog = `
SELECT format('column %s.%s %s%s%s', c.relname, a.attname, format_type(a.atttypid, a.atttypmod),
	CASE WHEN a.attnotnull THEN ' not null' ELSE '' END, ' default ' || pg_get_expr(d.adbin, d.adrelid))
FROM pg_attribute a
JOIN pg_class c ON c.oid = a.attrelid
LEFT JOIN pg_attrdef d ON d.adrelid = a.attrelid AND d.adnum = a.attnum
WHERE c.relnamespace = current_schema()::regnamespace AND c.relkind = 'r' AND c.relname LIKE 'auth\_%'
	AND a.attnum > 0 AND NOT a.attisdropped
UNION ALL
SELECT format('constraint %s.%s %s', c.relname, k.conname, pg_get_constraintdef(k.oid))
FROM pg_constraint k
JOIN pg_class c ON c.oid = k.conrelid
WHERE c.relnamespace = current_schema()::regnamespace AND c.relname LIKE 'auth\_%'
UNION ALL
SELECT format('index %s', pg_get_indexdef(i.indexrelid))
FROM pg_index i
JOIN pg_class c ON c.oid = i.indrelid
WHERE c.relnamespace = current_schema()::regnamespace AND c.relname LIKE 'auth\_%'
ORDER BY 1`

func TestSchema(t *testing.T) {
	t.Setenv("PGDATABASE", "pgtest")
	old, fresh := pgtest.Open(t), pgtest.Open(t)
	if err := migrate.Command(t.Context(), old, []migrate.Track{Track}, []string{"up"}, io.Discard); err != nil {
		t.Fatal(err)
	}
	run := func(db *gorm.DB) []string {
		t.Helper()
		r := &recorder{Interface: logger.Discard}
		if err := migrate.Up(t.Context(), db.Session(&gorm.Session{Logger: r}), []migrate.Schema{Schema}, io.Discard); err != nil {
			t.Fatal(err)
		}
		return r.ddl
	}
	lines := func(db *gorm.DB) []string {
		t.Helper()
		var lines []string
		if err := db.Raw(catalog).Scan(&lines).Error; err != nil {
			t.Fatal(err)
		}
		return lines
	}
	same := func(label string) {
		t.Helper()
		if want, got := lines(old), lines(fresh); !slices.Equal(got, want) {
			t.Errorf("%s: Schema catalog differs from Track\nTrack:\n%s\nSchema:\n%s", label, strings.Join(want, "\n"), strings.Join(got, "\n"))
		}
	}

	if len(run(fresh)) == 0 {
		t.Fatal("Run created nothing on an empty database")
	}
	same("fresh")
	if ddl := run(old); len(ddl) != 1 || !regexp.MustCompile(`^CREATE TABLE "hgmx_init"`).MatchString(ddl[0]) {
		t.Errorf("Run on the Track schema issued %q, want only the ledger table", ddl)
	}
	same("upgraded")
	for name, db := range map[string]*gorm.DB{"old": old, "fresh": fresh} {
		if ddl := run(db); len(ddl) > 0 {
			t.Errorf("rerun on %s issued %q", name, ddl)
		}
	}
}
