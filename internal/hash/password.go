package hash

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPassLen = 8
	maxPassLen = 72
)

var (
	allowedChars = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*]+$`)
	hasUpper     = regexp.MustCompile(`[A-Z]`)
	hasDigit     = regexp.MustCompile(`[0-9]`)
	hasSpecial   = regexp.MustCompile(`[!@#$%^&*]`)
)

var ErrInvalidPassword = errors.New("invalid password")

func ValidatePassword(pass string) (string, error) {
	if len(pass) < minPassLen || len(pass) > maxPassLen {
		return "", ErrInvalidPassword
	}
	if !allowedChars.MatchString(pass) || !hasUpper.MatchString(pass) ||
		!hasDigit.MatchString(pass) || !hasSpecial.MatchString(pass) {
		return "", ErrInvalidPassword
	}
	return pass, nil
}

func HashPassword(pass string) (string, error) {
	pass, err := ValidatePassword(pass)
	if err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			return "", ErrInvalidPassword
		}
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
