package pg

import (
	"log/slog"

	"github.com/cogikyo/hgmx/pg/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound   = gorm.ErrRecordNotFound
	ErrDuplicate  = gorm.ErrDuplicatedKey
	ErrForeignKey = gorm.ErrForeignKeyViolated
)

func Open(log *slog.Logger) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{}), config.Gorm(log))
}
