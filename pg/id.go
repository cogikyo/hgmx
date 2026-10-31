package pg

import (
	"database/sql/driver"
	"fmt"
	"uuid"
)

// ID stores a UUID as a Postgres uuid; its zero value maps to SQL NULL.
type ID struct{ uuid.UUID }

// NewID makes a UUIDv7.
func NewID() ID { return ID{uuid.NewV7()} }

// ParseID parses a UUID string.
func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	return ID{u}, err
}

// Scan accepts a UUID string or SQL NULL, which clears id.
func (id *ID) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*id = ID{}
		return nil
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = ID{u}
		return nil
	}
	return fmt.Errorf("pg: cannot scan %T into ID", src)
}

// Value returns a UUID string, or SQL NULL for the zero ID.
func (id ID) Value() (driver.Value, error) {
	if id == (ID{}) {
		return nil, nil
	}
	return id.String(), nil
}

// GormDataType selects Postgres uuid storage.
func (ID) GormDataType() string { return "uuid" }
