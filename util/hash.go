package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"hash"
)

type Hash []byte

func (h *Hash) Scan(val interface{}) error {
	if b, ok := val.([]byte); ok {
		if result, err := hex.DecodeString(string(b)); err != nil {
			return err
		} else {
			*h = result
		}
	}
	return nil
}

func (h Hash) Value() (driver.Value, error) {
	return h.String(), nil
}

func (h Hash) String() string {
	return hex.EncodeToString(h)
}

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	h.Scan(data)
	return nil
}

func NewHash(h hash.Hash, items ...string) Hash {
	for _, s := range items {
		h.Write([]byte(s))
	}

	return h.Sum(nil)
}

func NewHashString(h hash.Hash, items ...string) string {
	return NewHash(h, items...).String()
}

func MD5(items ...string) Hash {
	return NewHash(md5.New(), items...)
}

func MD5String(items ...string) string {
	return NewHashString(md5.New(), items...)
}

func SHA1(items ...string) Hash {
	return NewHash(sha1.New(), items...)
}

func SHA1String(items ...string) string {
	return NewHashString(sha1.New(), items...)
}

func SHA256(items ...string) Hash {
	return NewHash(sha256.New(), items...)
}

func SHA256String(items ...string) string {
	return NewHashString(sha256.New(), items...)
}
