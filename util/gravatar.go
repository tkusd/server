package util

func Gravatar(email string) string {
	return "https://www.gravatar.com/avatar/" + MD5(email).String()
}
