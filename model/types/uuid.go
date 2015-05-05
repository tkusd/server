package types

import (
	"database/sql/driver"

	"code.google.com/p/go-uuid/uuid"
)

// UUID inherits uuid.UUID.
type UUID struct {
	uuid.UUID
}

// Scan implements the sql.Scanner interface.
func (uid *UUID) Scan(val interface{}) error {
	if b, ok := val.([]byte); ok {
		uid.UUID = uuid.Parse(string(b))
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (uid UUID) Value() (driver.Value, error) {
	if uid.IsEmpty() {
		return nil, nil
	}

	return uid.UUID.String(), nil
}

// Equal checks the equality of two UUIDs.
func (uid UUID) Equal(a UUID) bool {
	return uuid.Equal(uid.UUID, a.UUID)
}

// IsEmpty check whether the UUID is empty or not.
func (uid UUID) IsEmpty() bool {
	return len(uid.UUID) == 0
}

// NewRandomUUID returns a random UUID (Version 4).
func NewRandomUUID() UUID {
	return UUID{uuid.NewRandom()}
}

// ParseUUID parses the string and returns a UUID.
func ParseUUID(id string) UUID {
	return UUID{uuid.Parse(id)}
}
