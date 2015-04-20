package util

import "code.google.com/p/go.crypto/bcrypt"

func GenerateBcryptHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareBcryptHash(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
