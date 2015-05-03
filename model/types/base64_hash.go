package types

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

// Base64Hash inherits from Hash, but returns Base64 string instead.
type Base64Hash struct {
	Hash
}

// Value implements the driver.Valuer interface.
func (h Base64Hash) Value() (driver.Value, error) {
	return h.HexString(), nil
}

// String implements the fmt.Stringer interface. It returns a string with base64 encoding.
func (h Base64Hash) String() string {
	return h.Base64String()
}

// HexString returns a string with hex encoding.
func (h Base64Hash) HexString() string {
	return hex.EncodeToString(h.Hash)
}

// Base64String returns a string with base64 encoding.
func (h Base64Hash) Base64String() string {
	return base64.URLEncoding.EncodeToString(h.Hash)
}

// IsEmpty returns true when the hash is empty.
func (h Base64Hash) IsEmpty() bool {
	return len(h.Hash) == 0
}

func (h Base64Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

func (h Base64Hash) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalJSON implements the json.Marshaler interface.
func (h *Base64Hash) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	result, err := DecodeBase64(str)

	if err != nil {
		return err
	}

	*h = *result

	return nil
}

// DecodeBase64 decodes the input string and returns a hash.
func DecodeBase64(key string) (*Base64Hash, error) {
	result, err := base64.URLEncoding.DecodeString(key)

	if err != nil {
		return nil, err
	}

	return &Base64Hash{result}, nil
}
