package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"hash"
)

// Hash equals to a slice of byte. But it can deal with the database properly.
type Hash []byte

// Scan implements the sql.Scanner interface.
func (h *Hash) Scan(val interface{}) error {
	if b, ok := val.([]byte); ok {
		result, err := hex.DecodeString(string(b))

		if err != nil {
			return err
		}

		*h = result
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (h Hash) Value() (driver.Value, error) {
	return h.String(), nil
}

// String implements the fmt.Stringer interface. It returns a string with hex encoding.
func (h Hash) String() string {
	return hex.EncodeToString(h)
}

// MarshalJSON implements json.Marshaler interface.
func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

// MarshalText implements encoding.TextMarshaler interface.
func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *Hash) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	result, err := hex.DecodeString(str)

	if err != nil {
		return err
	}

	*h = result

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (h *Hash) UnmarshalText(data []byte) error {
	return h.UnmarshalJSON(data)
}

// NewHash creates a new hash instance.
func NewHash(h hash.Hash, items ...string) Hash {
	for _, s := range items {
		h.Write([]byte(s))
	}

	return h.Sum(nil)
}

// MD5 creates a new MD5 hash.
func MD5(items ...string) Hash {
	return NewHash(md5.New(), items...)
}

// SHA1 creates a new SHA1 hash.
func SHA1(items ...string) Hash {
	return NewHash(sha1.New(), items...)
}

// SHA256 creates a new SHA256 hash.
func SHA256(items ...string) Hash {
	return NewHash(sha256.New(), items...)
}
