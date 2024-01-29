package env

import "regexp"

const (
	EmailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	UserNameRegex = `^[a-zA-Z ]+`
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(EmailRegex)
	return emailRegex.MatchString(email)
}

func IsValidName(name string) bool {
	nameRegex := regexp.MustCompile(UserNameRegex)
	return nameRegex.MatchString(name)
}
