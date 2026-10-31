package pg

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Page struct{ Offset, Limit int }

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

type ListResult[T any] struct {
	Count int64
	Items []T
}

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

func Like(column, q string) func(*gorm.DB) *gorm.DB {
	q = strings.TrimSpace(q)
	return func(db *gorm.DB) *gorm.DB {
		if q == "" {
			return db
		}
		return db.Where("? ILIKE ?", clause.Column{Name: column}, "%"+likeEscaper.Replace(q)+"%")
	}
}

func In[T any](column string, vs []T) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(vs) == 0 {
			return db.Where("false")
		}
		return db.Where("? IN ?", clause.Column{Name: column}, vs)
	}
}

func Sort(allowed map[string]string, by string, desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		column, ok := allowed[by]
		if !ok {
			return db
		}
		return db.Order(clause.OrderByColumn{Column: clause.Column{Name: column}, Desc: desc})
	}
}
