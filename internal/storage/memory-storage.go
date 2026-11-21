package storage

import (
	"sync"

	"github.com/kimkachailo/auth-service/internal/domain"
)

type Database struct {
	storage   map[int]domain.User
	counterID int
	mu        sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		storage: make(map[int]domain.User),
	}
}

func (d *Database) AddUser(emailAddress string, hash string, timeStamp string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, v := range d.storage {
		if v.Email == emailAddress {
			return domain.ErrEmailAlreadyExists
		}
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

func (d *Database) SearchID(ID int) (domain.User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	v, ok := d.storage[ID]
	if ok {
		return v, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (d *Database) SearchEmail(emailAddress string) (domain.User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for k, v := range d.storage {
		if v.Email == emailAddress {
			return d.storage[k], nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}
