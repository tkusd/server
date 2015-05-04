package types

import (
	"database/sql/driver"
	"time"
)

// Time inherits from time.Time. But prints ISO 8601 format time in json.
type Time struct {
	time.Time
}

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(data interface{}) error {
	if val, ok := data.(time.Time); ok {
		t.Time = val
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return t.Time.Format(time.RFC3339Nano), nil
}

// MarshalJSON implements json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.ISOTime() + `"`), nil
}

// MarshalText implements encoding.TextMarshaler interface.
func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.ISOTime()), nil
}

// ISOTime returns time in ISO 8601 time format.
func (t Time) ISOTime() string {
	return ISOTime(t.Time)
}

// ISOTime formats a time in ISO 8601 time format.
func ISOTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ParseISOTime parses a ISO 8601 format string to time.
func ParseISOTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

// Now returns current time.
func Now() Time {
	return Time{time.Now().UTC()}
}
