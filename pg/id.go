package pg

import (
	"database/sql/driver"
	"fmt"
	"uuid"
)

type ID struct{ uuid.UUID }

func NewID() ID { return ID{uuid.NewV7()} }

func ParseID(s string) (ID, error) {
	u, err := uuid.Parse(s)
	return ID{u}, err
}

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

func (id ID) Value() (driver.Value, error) {
	if id == (ID{}) {
		return nil, nil
	}
	return id.String(), nil
}

func (ID) GormDataType() string { return "uuid" }
