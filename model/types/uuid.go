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
	if !uid.Valid() {
		return nil, nil
	}

	return uid.UUID.String(), nil
}

/*
func (uid UUID) MarshalJSON() ([]byte, error) {
	if !uid.Valid() {
		return []byte("''"), nil
	}

	return []byte(`"` + uid.String() + `"`), nil
}*/

// Equal checks the equality of two UUIDs.
func (uid UUID) Equal(a UUID) bool {
	return uuid.Equal(uid.UUID, a.UUID)
}

// Valid returns true if UUID is valid.
func (uid UUID) Valid() bool {
	return uid.UUID.Variant() != uuid.Invalid
}

// NewRandomUUID returns a random UUID (Version 4).
func NewRandomUUID() UUID {
	return UUID{uuid.NewRandom()}
}

// ParseUUID parses the string and returns a UUID.
func ParseUUID(id string) UUID {
	return UUID{uuid.Parse(id)}
}
