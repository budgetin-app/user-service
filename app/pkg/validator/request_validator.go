package validator

import "regexp"

const (
	UsernameRegex = `^[a-zA-Z0-9]*$`
	EmailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	UserNameRegex = `^[a-zA-Z ]+`
)

func IsValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(UsernameRegex)
	return usernameRegex.MatchString(username)
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(EmailRegex)
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	// Check if the password is at least 8 characters long
	if len(password) < 8 {
		return false
	}

	// Check if the password contains at least one digit
	digitRegex := regexp.MustCompile(`[0-9]`)
	if !digitRegex.MatchString(password) {
		return false
	}

	// Check if the password contains at least one special character
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()-_+=]`)
	return specialCharRegex.MatchString(password)
}

func IsValidName(name string) bool {
	nameRegex := regexp.MustCompile(UserNameRegex)
	return nameRegex.MatchString(name)
}
