package config

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode/utf8"

	"gorm.io/gorm/schema"
)

// Match Postgres defaults up to 63 bytes; hash longer names, unlike Postgres.
type naming struct{ schema.NamingStrategy }

func (n naming) CheckerName(table, column string) string {
	return name(table, column, "check")
}

func (n naming) IndexName(table, column string) string {
	return name(table, n.ColumnName(table, column), "idx")
}

func (n naming) UniqueName(table, column string) string {
	return name(table, n.ColumnName(table, column), "key")
}

func (n naming) RelationshipFKName(rel schema.Relationship) string {
	table := rel.Schema.Table
	var columns []string
	for _, ref := range rel.References {
		if ref.PrimaryKey != nil && (rel.JoinTable == nil || ref.OwnPrimaryKey) {
			table = ref.ForeignKey.Schema.Table
			columns = append(columns, ref.ForeignKey.DBName)
		}
	}
	return name(table, strings.Join(columns, "_"), "fkey")
}

func name(table, columns, label string) string {
	full := table + "_" + columns + "_" + label
	if len(full) <= 63 {
		return full
	}
	sum := sha256.Sum256([]byte(full))
	suffix := fmt.Sprintf("_%x_%s", sum[:4], label)
	return clip(full, 63-len(suffix)) + suffix
}

func clip(s string, n int) string {
	for n < len(s) && !utf8.RuneStart(s[n]) {
		n--
	}
	return s[:n]
}
