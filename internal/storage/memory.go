package storage

import (
	"errors"
	"fmt"

	"github.com/kimkachailo/auth-service/internal/domain"
)

type Database struct {
	storage   map[int]domain.User
	counterID int
}

func NewDatabase() *Database {
	return &Database{
		storage: make(map[int]domain.User),
	}
}

func (d *Database) AddUser(emailAddress string, hash string, timeStamp string) error {
	_, err := d.SearchEmail(emailAddress)
	if err == nil {
		return fmt.Errorf("in AddUser: %w", domain.ErrEmailAlreadyExists)
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return err
	}
	d.counterID++
	d.storage[d.counterID] = domain.User{
		ID:           d.counterID,
		Email:        emailAddress,
		PasswordHash: hash,
		CreatedAt:    timeStamp,
	}
	return nil
}

func (d Database) SearchID(ID int) (domain.User, error) {
	v, ok := d.storage[ID]
	if ok {
		return v, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (d Database) SearchEmail(emailAddress string) (domain.User, error) {
	for k, v := range d.storage {
		if v.Email == emailAddress {
			return d.storage[k], nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}
