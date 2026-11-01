// Package migrate runs Go-defined Postgres schema steps through goose.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"gorm.io/gorm"
)

// Step is a schema change; versions run 1, 2, 3, ... per track in list order.
// Steps run in transactions with a 5-second lock_timeout unless NoTx is set.
type Step struct {
	Version int64
	Name    string
	// Up is required; a nil Down makes rollback fail without changing history.
	Up, Down func(ctx context.Context, db *gorm.DB) error
	// NoTx disables the transaction; the step must tolerate a partial rerun.
	NoTx bool
}

// Track groups steps under an independent migration history.
// An empty Name uses goose_db_version; other names use hgmx_<name>_migrations.
// Names must match [a-z][a-z0-9_]{0,46} unless empty.
type Track struct {
	Name  string
	Steps []Step
}

const usage = "usage: migrate up | down [track] | status"

var trackName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,46}$`)

// Command runs up in track order, down by one step, or status for every track.
// Down defaults to the unnamed site track; pass a name to target another.
// With steps, PGDATABASE is required and one session lock covers the command.
// The pool needs two open connections, or three when running NoTx steps.
func Command(ctx context.Context, db *gorm.DB, tracks []Track, args []string, out io.Writer) (retErr error) {
	if len(args) == 0 || len(args) > 2 || (args[0] != "up" && args[0] != "down" && args[0] != "status") || (len(args) == 2 && args[0] != "down") {
		return errors.New(usage)
	}
	target := ""
	if len(args) == 2 {
		target = args[1]
	}
	names := make(map[string]bool, len(tracks))
	count := 0
	noTx := false
	for _, track := range tracks {
		if track.Name != "" && !trackName.MatchString(track.Name) {
			return fmt.Errorf("invalid track name %q; want [a-z][a-z0-9_]{0,46}", track.Name)
		}
		if names[track.Name] {
			return fmt.Errorf("duplicate track name %q", track.Name)
		}
		names[track.Name] = true
		for i, s := range track.Steps {
			if s.Version != int64(i+1) {
				return fmt.Errorf("track %q: step %d (%s) has version %d; versions run 1, 2, 3, … in list order", track.Name, i+1, s.Name, s.Version)
			}
			if s.Up == nil {
				return fmt.Errorf("track %q: step %d (%s) has no Up", track.Name, s.Version, s.Name)
			}
			if s.NoTx && (args[0] == "up" || (args[0] == "down" && track.Name == target)) {
				noTx = true
			}
		}
		count += len(track.Steps)
	}
	if len(args) == 2 && !names[target] {
		return fmt.Errorf("unknown track %q", target)
	}
	if count == 0 {
		if args[0] == "status" && len(tracks) > 0 {
			for _, track := range tracks {
				fmt.Fprintf(out, "%s\nno migrations\n", trackLabel(track.Name))
			}
		} else if args[0] == "up" || args[0] == "status" {
			fmt.Fprintln(out, "no migrations")
		} else {
			fmt.Fprintln(out, "nothing to roll back")
		}
		return nil
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

	minConns := 2
	if noTx {
		minConns = 3
	}
	if max := conn.Stats().MaxOpenConnections; max > 0 && max < minConns {
		return fmt.Errorf("migrations require at least %d open connections for the session lock and migration work", minConns)
	}
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return err
	}
	locked, err := conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer locked.Close()
	if err := locker.SessionLock(ctx, locked); err != nil {
		return err
	}
	defer func() {
		retErr = errors.Join(retErr, locker.SessionUnlock(context.WithoutCancel(ctx), locked))
	}()

	if args[0] == "down" && !names[target] {
		fmt.Fprintln(out, "nothing to roll back")
		return nil
	}
	for _, track := range tracks {
		if args[0] == "down" && track.Name != target {
			continue
		}
		fmt.Fprintln(out, trackLabel(track.Name))
		if len(track.Steps) == 0 {
			if args[0] == "down" {
				fmt.Fprintln(out, "nothing to roll back")
			} else {
				fmt.Fprintln(out, "no migrations")
			}
			continue
		}
		p, err := provider(db, conn, track)
		if err != nil {
			return err
		}
		switch args[0] {
		case "up":
			rs, err := p.Up(ctx)
			report(out, track.Steps, rs, err)
			if err != nil {
				return err
			}
			if len(rs) == 0 {
				fmt.Fprintln(out, "no pending migrations")
			}
		case "down":
			r, err := p.Down(ctx)
			if errors.Is(err, goose.ErrNoNextVersion) {
				fmt.Fprintln(out, "nothing to roll back")
				continue
			}
			if err != nil {
				return err
			}
			report(out, track.Steps, []*goose.MigrationResult{r}, nil)
		case "status":
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
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", s.Source.Version, s.State, applied, track.Steps[s.Source.Version-1].Name)
			}
			if err := w.Flush(); err != nil {
				return err
			}
		}
	}
	return nil
}

func report(out io.Writer, steps []Step, rs []*goose.MigrationResult, err error) {
	if pe, ok := errors.AsType[*goose.PartialError](err); ok {
		rs = pe.Applied
	}
	for _, r := range rs {
		fmt.Fprintf(out, "%s %d %s (%s)\n", r.Direction, r.Source.Version, steps[r.Source.Version-1].Name, r.Duration.Round(time.Millisecond))
	}
}

func trackLabel(name string) string {
	if name == "" {
		return "track site"
	}
	return "track " + name
}

func provider(db *gorm.DB, conn *sql.DB, track Track) (*goose.Provider, error) {
	table := goose.DefaultTablename
	if track.Name != "" {
		table = "hgmx_" + track.Name + "_migrations"
	}
	steps := track.Steps
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
		goose.WithTableName(table),
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
