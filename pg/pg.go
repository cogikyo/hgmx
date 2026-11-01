// Package pg supplies Postgres connections, column types and query scopes.
package pg

import (
	"context"
	"log/slog"

	"github.com/cogikyo/hgmx/pg/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// ErrNotFound is the no-record outcome, tested with errors.Is.
	ErrNotFound = gorm.ErrRecordNotFound
	// ErrDuplicate is a unique-constraint violation, tested with errors.Is.
	ErrDuplicate = gorm.ErrDuplicatedKey
	// ErrForeignKey is a foreign-key violation, tested with errors.Is.
	ErrForeignKey = gorm.ErrForeignKeyViolated
)

// Open reads libpq-compatible PG* settings and opens a pool without pinging.
// The caller must [Close] the pool; SQL logs use placeholders, not bound values.
func Open(log *slog.Logger) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{}), config.Gorm(log))
}

// Close closes db's pool and returns ctx.Err() if ctx ends first.
// Pool closure continues after cancellation; drain requests before calling.
func Close(ctx context.Context, db *gorm.DB) error {
	pool, err := db.DB()
	if err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- pool.Close() }()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
