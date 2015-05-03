package util

import "code.google.com/p/go.crypto/bcrypt"

// GenerateBcryptHash generates a bcrypt hash.
func GenerateBcryptHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CompareBcryptHash compares bcrypt hash and the password.
func CompareBcryptHash(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
