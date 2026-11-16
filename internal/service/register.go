package service

import (
	"auth-service/internal/hash"
	"time"
)

type UserAdder interface {
	AddUser(emailAddress string, hash string, timeStamp string) error
}

func Register(emailAddress string, password string, db UserAdder) error {
	pHash, err := hash.Password(password)
	if err != nil {
		return err
	}
	timeStamp := time.Now().UTC().Truncate(time.Second).String()
	if err := db.AddUser(emailAddress, pHash, timeStamp); err != nil {
		return err
	}
	return nil
}
