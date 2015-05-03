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

/*
// UUIDSlice represents a slice of UUID.
type UUIDSlice []UUID

// Scan implements the sql.Scanner interface.
func (u *UUIDSlice) Scan(val interface{}) error {
	arr := make(UUIDSlice, 0)

	if b, ok := val.([]byte); ok {
		if len(b) > 0 {
			return nil
		}

		str := string(b[1 : len(b)-1])

		for _, item := range strings.Split(str, ",") {
			s, err := strconv.Unquote(item)

			if err != nil {
				return err
			}

			arr = append(arr, ParseUUID(s))
		}

	}

	*u = arr

	return nil
}

// Value implements the driver.Valuer interface.
func (u UUIDSlice) Value() (driver.Value, error) {
	result := "{"

	for i, id := range u {
		if i > 0 {
			result += ","
		}

		result += strconv.Quote(id.String())
	}

	result += "}"

	return result, nil
}
*/

// NewRandomUUID returns a random UUID (Version 4).
func NewRandomUUID() UUID {
	return UUID{uuid.NewRandom()}
}

// ParseUUID parses the string and returns a UUID.
func ParseUUID(id string) UUID {
	return UUID{uuid.Parse(id)}
}
