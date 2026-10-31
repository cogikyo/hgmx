package pg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSON[T any] struct{ V T }

func (j *JSON[T]) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case nil:
		*j = JSON[T]{}
		return nil
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("pg: cannot scan %T into JSON", src)
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	j.V = v
	return nil
}

func (j JSON[T]) Value() (driver.Value, error) {
	b, err := json.Marshal(j.V)
	return string(b), err
}

func (JSON[T]) GormDBDataType(*gorm.DB, *schema.Field) string { return "jsonb" }

func (j JSON[T]) MarshalJSON() ([]byte, error) { return json.Marshal(j.V) }

func (j *JSON[T]) UnmarshalJSON(b []byte) error { return json.Unmarshal(b, &j.V) }
