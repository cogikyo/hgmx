package config

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Gorm(log *slog.Logger) *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		TranslateError:       true,
		NowFunc:              func() time.Time { return time.Now().UTC() },
		Logger: logger.NewSlogLogger(log, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
	}
}
