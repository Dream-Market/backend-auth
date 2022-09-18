package validation

import (
	"net/mail"
	"unicode"
)

func ValidateEmail(login string) (email string, ok bool) {
	address, err := mail.ParseAddress(login)
	if err != nil {
		return "", false
	}
	return address.Address, true
}

func ValidatePhone(login string) (phone string, ok bool) {
	if len(login) != 11 {
		return "", false
	}

	for _, ch := range login {
		if !unicode.IsDigit(ch) {
			return phone, false
		}
	}
	return login, true
}

func ValidatePassword(password string) (ok bool) {
	if len(password) < 8 {
		return false
	}

	var upper, lower, digit bool

	for _, ch := range password {
		if unicode.IsUpper(ch) {
			upper = true
		} else if unicode.IsLower(ch) {
			lower = true
		} else if unicode.IsDigit(ch) {
			digit = true
		}
	}

	return upper && lower && digit
}
