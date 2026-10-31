// Package migrate runs Go-defined Postgres schema steps through goose.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"gorm.io/gorm"
)

// Step is a schema change; versions must run 1, 2, 3, ... in list order.
// Steps run in transactions with a 5-second lock_timeout unless NoTx is set.
type Step struct {
	Version int64
	Name    string
	// Up is required; a nil Down makes rollback fail without changing history.
	Up, Down func(ctx context.Context, db *gorm.DB) error
	// NoTx disables the transaction; the step must tolerate a partial rerun.
	NoTx bool
}

const usage = "usage: migrate up | down | status"

// Command runs up, down (one step) or status and writes to out.
// It requires PGDATABASE and pings db, even when steps is empty.
// With steps registered, all commands take the session lock, including status.
func Command(ctx context.Context, db *gorm.DB, steps []Step, args []string, out io.Writer) error {
	if len(args) != 1 || (args[0] != "up" && args[0] != "down" && args[0] != "status") {
		return errors.New(usage)
	}
	for i, s := range steps {
		if s.Version != int64(i+1) {
			return fmt.Errorf("step %d (%s) has version %d; versions run 1, 2, 3, … in list order", i+1, s.Name, s.Version)
		}
		if s.Up == nil {
			return fmt.Errorf("step %d (%s) has no Up", s.Version, s.Name)
		}
	}
	if os.Getenv("PGDATABASE") == "" {
		return errors.New("PGDATABASE is unset; refusing to migrate the default database")
	}
	conn, err := db.DB()
	if err != nil {
		return err
	}
	if err := conn.PingContext(ctx); err != nil {
		return err
	}

	if len(steps) == 0 {
		if args[0] == "up" || args[0] == "status" {
			fmt.Fprintln(out, "no migrations")
		} else {
			fmt.Fprintln(out, "nothing to roll back")
		}
		return nil
	}
	p, err := provider(db, conn, steps)
	if err != nil {
		return err
	}
	switch args[0] {
	case "up":
		rs, err := p.Up(ctx)
		report(out, steps, rs, err)
		if err == nil && len(rs) == 0 {
			fmt.Fprintln(out, "no pending migrations")
		}
		return err
	case "down":
		r, err := p.Down(ctx)
		if errors.Is(err, goose.ErrNoNextVersion) {
			fmt.Fprintln(out, "nothing to roll back")
			return nil
		}
		if err != nil {
			return err
		}
		report(out, steps, []*goose.MigrationResult{r}, nil)
		return nil
	}
	ss, err := p.Status(ctx)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, s := range ss {
		applied := "-"
		if s.State == goose.StateApplied {
			applied = s.AppliedAt.UTC().Format(time.DateTime)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", s.Source.Version, s.State, applied, steps[s.Source.Version-1].Name)
	}
	return w.Flush()
}

func report(out io.Writer, steps []Step, rs []*goose.MigrationResult, err error) {
	if pe, ok := errors.AsType[*goose.PartialError](err); ok {
		rs = pe.Applied
	}
	for _, r := range rs {
		fmt.Fprintf(out, "%s %d %s (%s)\n", r.Direction, r.Source.Version, steps[r.Source.Version-1].Name, r.Duration.Round(time.Millisecond))
	}
}

func provider(db *gorm.DB, conn *sql.DB, steps []Step) (*goose.Provider, error) {
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return nil, err
	}
	ms := make([]*goose.Migration, len(steps))
	for i, s := range steps {
		down := s.Down
		if down == nil {
			down = func(context.Context, *gorm.DB) error {
				return fmt.Errorf("step %d (%s) has no Down", s.Version, s.Name)
			}
		}
		ms[i] = goose.NewGoMigration(s.Version, run(db, s.NoTx, s.Up), run(db, s.NoTx, down))
	}
	return goose.NewProvider(goose.DialectPostgres, conn, nil,
		goose.WithSessionLocker(locker),
		goose.WithDisableGlobalRegistry(true),
		goose.WithGoMigrations(ms...),
	)
}

func run(db *gorm.DB, noTx bool, f func(context.Context, *gorm.DB) error) *goose.GoFunc {
	if noTx {
		return &goose.GoFunc{RunDB: func(ctx context.Context, _ *sql.DB) error {
			return f(ctx, db.WithContext(ctx))
		}}
	}
	return &goose.GoFunc{RunTx: func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, "SET LOCAL lock_timeout = '5s'"); err != nil {
			return err
		}
		g := db.Session(&gorm.Session{Context: ctx, NewDB: true})
		g.Statement.ConnPool = tx
		return f(ctx, g)
	}}
}
