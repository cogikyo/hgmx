package migrate

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"slices"

	"github.com/pressly/goose/v3/lock"
	"gorm.io/gorm"

	"github.com/cogikyo/hgmx/cli"
	"github.com/cogikyo/hgmx/pg"
)

// Schema registers a package's target models, enums and optional init SQL.
type Schema struct {
	Name   string
	Enums  []Enum
	Models []any
	Init   fs.FS // Embed the init directory, or leave nil when no SQL is needed.
}

// Enum creates or extends an enum before model and init changes run.
// [github.com/cogikyo/hgmx/pg.Enum] implements Enum.
type Enum interface {
	Create(tx *gorm.DB) error
}

type ledger struct {
	Schema string `gorm:"primaryKey"`
	File   string `gorm:"primaryKey"`
	Sha256 string `gorm:"not null"`
}

func (ledger) TableName() string { return "hgmx_init" }

type script struct {
	file string
	sql  []byte
}

type phases struct{ pre, post []script }

// Run returns a migrate command for [cli.Dispatch] that runs [Up] on "up".
// It opens and closes a pool from PG* env via [pg.Open] and rejects other args.
func Run(schemas ...Schema) cli.Command {
	return func(ctx context.Context, args []string) (err error) {
		if len(args) > 0 && (args[0] == "down" || args[0] == "status") {
			return cli.Usage("migrate %s is unsupported: models and init files define the schema, and up applies them", args[0])
		}
		if len(args) != 1 || args[0] != "up" {
			return cli.Usage("usage: migrate up")
		}
		db, err := pg.Open(slog.Default())
		if err != nil {
			return err
		}
		defer func() { err = errors.Join(err, pg.Close(context.WithoutCancel(ctx), db)) }()
		return Up(ctx, db, schemas, os.Stdout)
	}
}

// Up applies schemas to db in slice order; PGDATABASE must be set.
// One connection holds the migration lock across two transactions.
// The first commits all enums, which survive failure of the second.
// The second runs pre SQL, AutoMigrate and post SQL, then rewrites hgmx_init.
// Init files run on every up and must be idempotent and transactional.
func Up(ctx context.Context, db *gorm.DB, schemas []Schema, out io.Writer) (err error) {
	files := make([]phases, len(schemas))
	names := make(map[string]bool, len(schemas))
	for i, s := range schemas {
		if s.Name == "" {
			return fmt.Errorf("schema %d has no name", i+1)
		}
		if names[s.Name] {
			return fmt.Errorf("duplicate schema %q", s.Name)
		}
		names[s.Name] = true
		if files[i], err = read(s.Init); err != nil {
			return fmt.Errorf("schema %s: %w", s.Name, err)
		}
	}
	if os.Getenv("PGDATABASE") == "" {
		return errors.New("PGDATABASE is unset; refusing to migrate the default database")
	}

	err = db.Session(&gorm.Session{Context: ctx, NewDB: true}).Connection(func(conn *gorm.DB) (err error) {
		// Share legacy Command's lock.
		if err := conn.Exec("SELECT pg_advisory_lock(?)", lock.DefaultLockID).Error; err != nil {
			return err
		}
		defer func() {
			unlock := conn.WithContext(context.WithoutCancel(ctx)).Exec("SELECT pg_advisory_unlock(?)", lock.DefaultLockID).Error
			err = errors.Join(err, unlock)
		}()

		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("SET LOCAL lock_timeout = '5s'").Error; err != nil {
				return err
			}
			for _, s := range schemas {
				for _, e := range s.Enums {
					if err := e.Create(tx); err != nil {
						return fmt.Errorf("schema %s: enum: %w", s.Name, err)
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("SET LOCAL lock_timeout = '5s'").Error; err != nil {
				return err
			}
			var xid string
			if err := tx.Raw("SELECT pg_current_xact_id()::text").Scan(&xid).Error; err != nil {
				return err
			}
			exec := func(s Schema, scripts []script) error {
				for _, f := range scripts {
					if _, err := tx.Statement.ConnPool.ExecContext(ctx, string(f.sql)); err != nil {
						return fmt.Errorf("schema %s: %s: %w", s.Name, f.file, err)
					}
					// Detect SQL that ended the transaction instead of trusting Exec.
					var now string
					if err := tx.Raw("SELECT pg_current_xact_id()::text").Scan(&now).Error; err != nil {
						return fmt.Errorf("schema %s: %s: %w", s.Name, f.file, err)
					}
					if now != xid {
						return fmt.Errorf("schema %s: %s ended the migration transaction; init files must not COMMIT or ROLLBACK", s.Name, f.file)
					}
				}
				return nil
			}
			var rows []ledger
			for i, s := range schemas {
				if err := exec(s, files[i].pre); err != nil {
					return err
				}
				if err := tx.AutoMigrate(s.Models...); err != nil {
					return fmt.Errorf("schema %s: models: %w", s.Name, err)
				}
				if err := exec(s, files[i].post); err != nil {
					return err
				}
				for _, f := range slices.Concat(files[i].pre, files[i].post) {
					rows = append(rows, ledger{s.Name, f.file, fmt.Sprintf("%x", sha256.Sum256(f.sql))})
				}
			}
			if err := tx.AutoMigrate(&ledger{}); err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM hgmx_init").Error; err != nil {
				return err
			}
			if len(rows) == 0 {
				return nil
			}
			return tx.Create(&rows).Error
		})
	})
	if err != nil {
		return err
	}
	for i, s := range schemas {
		fmt.Fprintf(out, "up %s (%d models, %d init files)\n", s.Name, len(s.Models), len(files[i].pre)+len(files[i].post))
	}
	return nil
}

func read(fsys fs.FS) (phases, error) {
	var p phases
	if fsys == nil {
		return p, nil
	}
	root, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return p, err
	}
	for _, e := range root {
		if e.Name() != "init" || !e.IsDir() {
			return p, fmt.Errorf("%s is outside init/", e.Name())
		}
	}
	if fsys, err = fs.Sub(fsys, "init"); err != nil {
		return p, err
	}
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return p, err
	}
	for _, e := range entries {
		var scripts *[]script
		switch e.Name() {
		case "pre":
			scripts = &p.pre
		case "post":
			scripts = &p.post
		}
		if scripts == nil || !e.IsDir() {
			return p, fmt.Errorf("init holds %s; want only pre/ and post/ directories", e.Name())
		}
		files, err := fs.ReadDir(fsys, e.Name())
		if err != nil {
			return p, err
		}
		for _, f := range files {
			file := path.Join(e.Name(), f.Name())
			if !f.Type().IsRegular() || path.Ext(file) != ".sql" {
				return p, fmt.Errorf("init holds %s; want only .sql files", file)
			}
			sql, err := fs.ReadFile(fsys, file)
			if err != nil {
				return p, err
			}
			*scripts = append(*scripts, script{file, sql})
		}
	}
	return p, nil
}
