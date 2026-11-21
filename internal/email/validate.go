package email

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

const maxEmailLen = 254

var ErrInvalidEmail = errors.New("invalid email")

func Validate(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if len(email) > maxEmailLen {
		return "", ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidEmail, err)
	}
	if addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil

}
