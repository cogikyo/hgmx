// Package pg supplies Postgres connections, column types and query scopes.
package pg

import (
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
// The caller owns the pool; SQL logs use placeholders, not bound values.
func Open(log *slog.Logger) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{}), config.Gorm(log))
}
