package util

import (
	"database/sql/driver"

	"errors"

	"code.google.com/p/go-uuid/uuid"
)

type UUID struct {
	uuid.UUID
}

func (uid *UUID) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		uid.UUID = uuid.Parse(string(v))
		break

	case string:
		uid.UUID = uuid.Parse(v)
		break

	default:
		return errors.New("Incompatible type for UUID")
	}

	return nil
}

func (uid UUID) Value() (driver.Value, error) {
	return uid.UUID.String(), nil
}

func NewRandomUUID() UUID {
	return UUID{uuid.NewRandom()}
}

func ParseUUID(id string) UUID {
	return UUID{uuid.Parse(id)}
}
