package types

import (
	"database/sql/driver"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) Scan(data interface{}) error {
	if val, ok := data.(time.Time); ok {
		t.Time = val
	}

	return nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time.Format(time.RFC3339Nano), nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.ISOTime() + `"`), nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.ISOTime()), nil
}

func (t Time) ISOTime() string {
	return ISOTime(t.Time)
}

// ISOTime formats a time with ISO 8601 time format.
func ISOTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ParseISOTime parses a ISO 8601 format string to time.
func ParseISOTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

func Now() Time {
	return Time{time.Now().UTC()}
}
