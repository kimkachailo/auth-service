package service

import (
	"errors"
	"time"

	"github.com/kimkachailo/auth-service/internal/email"
	"github.com/kimkachailo/auth-service/internal/hash"
)

type UserAdder interface {
	AddUser(emailAddress string, hash string, timeStamp string) error
}

func SignUp(emailAddress string, password string, db UserAdder) error {
	validatedAddress, err := email.Validate(emailAddress)
	if err != nil {
		if errors.Is(err, email.ErrInvalidEmail) {
			return email.ErrInvalidEmail
		}
		return err
	}
	pHash, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	timeStamp := time.Now().UTC().Truncate(time.Second).String()
	if err := db.AddUser(validatedAddress, pHash, timeStamp); err != nil {
		return err
	}
	return nil
}
