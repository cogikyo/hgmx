package pg

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Page selects an offset and limit, normalized by Scope.
type Page struct{ Offset, Limit int }

// Scope floors Offset at zero, defaults Limit to 20 and caps it at 100.
// A non-positive Limit uses the default.
func (p Page) Scope(db *gorm.DB) *gorm.DB {
	limit := p.Limit
	switch {
	case limit <= 0:
		limit = 20
	case limit > 100:
		limit = 100
	}
	return db.Offset(max(p.Offset, 0)).Limit(limit)
}

// ListResult holds the query count before paging and the selected items.
type ListResult[T any] struct {
	Count int64
	Items []T
}

// List counts db's query before applying p to fetch Items.
// Count and Items come from separate queries.
func List[T any](db *gorm.DB, p Page) (ListResult[T], error) {
	var r ListResult[T]
	db = db.Session(&gorm.Session{})
	if err := db.Session(&gorm.Session{NewDB: true}).Table("(?) AS q", db.Model(new(T))).Count(&r.Count).Error; err != nil {
		return r, err
	}
	err := db.Scopes(p.Scope).Find(&r.Items).Error
	return r, err
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// Like matches a literal, case-insensitive substring of column.
// It trims whitespace; blank input adds no filter.
func Like(column, q string) func(*gorm.DB) *gorm.DB {
	q = strings.TrimSpace(q)
	return func(db *gorm.DB) *gorm.DB {
		if q == "" {
			return db
		}
		return db.Where("? ILIKE ?", clause.Column{Name: column}, "%"+likeEscaper.Replace(q)+"%")
	}
}

// In filters column to vs; an empty or nil slice matches nothing.
func In[T any](column string, vs []T) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(vs) == 0 {
			return db.Where("false")
		}
		return db.Where("? IN ?", clause.Column{Name: column}, vs)
	}
}

// Sort orders by an allowlisted key's column; an unknown key adds no order.
func Sort(allowed map[string]string, by string, desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		column, ok := allowed[by]
		if !ok {
			return db
		}
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: desc})
	}
}
