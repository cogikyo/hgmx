package pg

import (
	"slices"
	"strings"

	"gorm.io/gorm"
)

type Enum[T ~string] struct {
	Name   string
	Values []T
}

func (e Enum[T]) Valid(v T) bool { return slices.Contains(e.Values, v) }

func (e Enum[T]) Create(tx *gorm.DB) error {
	var exists bool
	err := tx.Raw(`SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ? AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = current_schema()))`, e.Name).Scan(&exists).Error
	if err != nil {
		return err
	}
	name := `"` + strings.ReplaceAll(e.Name, `"`, `""`) + `"`
	if !exists {
		labels := make([]string, len(e.Values))
		for i, v := range e.Values {
			labels[i] = literal(v)
		}
		return tx.Exec("CREATE TYPE " + name + " AS ENUM (" + strings.Join(labels, ", ") + ")").Error
	}
	for _, v := range e.Values {
		if err := tx.Exec("ALTER TYPE " + name + " ADD VALUE IF NOT EXISTS " + literal(v)).Error; err != nil {
			return err
		}
	}
	return nil
}

func literal[T ~string](v T) string {
	return "'" + strings.ReplaceAll(string(v), "'", "''") + "'"
}
