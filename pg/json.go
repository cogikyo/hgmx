package pg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JSON stores V as jsonb and marshals as V without a wrapper object.
// SQL NULL reads as zero T and saves as non-NULL JSON, including JSON null.
type JSON[T any] struct{ V T }

// Scan reads JSON bytes or text; SQL NULL clears V to the zero value of T.
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

// Value encodes V as JSON text, never SQL NULL.
func (j JSON[T]) Value() (driver.Value, error) {
	b, err := json.Marshal(j.V)
	return string(b), err
}

// GormDBDataType selects Postgres jsonb storage.
func (JSON[T]) GormDBDataType(*gorm.DB, *schema.Field) string { return "jsonb" }

// MarshalJSON encodes V directly.
func (j JSON[T]) MarshalJSON() ([]byte, error) { return json.Marshal(j.V) }

// UnmarshalJSON decodes directly into V.
func (j *JSON[T]) UnmarshalJSON(b []byte) error { return json.Unmarshal(b, &j.V) }
