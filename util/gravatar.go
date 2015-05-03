package util

import "github.com/tommy351/app-studio-server/model/types"

// Gravatar generates gravatar URL with email.
func Gravatar(email string) string {
	return "https://www.gravatar.com/avatar/" + types.MD5(email).String()
}
