package storage_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/kimkachailo/auth-service/internal/domain"
	"github.com/kimkachailo/auth-service/internal/storage"
)

func createUser(t *testing.T, email string, hash string, timeStamp string) *storage.Database {
	t.Helper()
	db := storage.NewDatabase()
	if err := db.AddUser(email, hash, timeStamp); err != nil {
		t.Fatalf("add user: got %v; want nil", err)
	}
	user, err := db.SearchEmail(email)
	if err != nil {
		t.Errorf("search email: got %v; want nil", err)
	}
	if user.Email != email {
		t.Fatalf("search email: got %s; want %s", user.Email, email)
	}
	return db
}

func TestStorage(t *testing.T) {
	t.Run("add and search", func(t *testing.T) {
		db := createUser(t, "user@example.com", "hash", "timeStamp")

		if _, err := db.SearchEmail("missing@example.com"); !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("search email: got %v; want %v", err, domain.ErrUserNotFound)
		}
		user, err := db.SearchID(1)
		if err != nil {
			t.Errorf("search id: got %v; want nil", err)
		}
		if user.ID != 1 {
			t.Fatalf("search id: got %d; want 1", user.ID)
		}
		if _, err := db.SearchID(1111); !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("search id: got %v; want %v", err, domain.ErrUserNotFound)
		}

	})
	t.Run("add duplicate", func(t *testing.T) {
		db := createUser(t, "user@example.com", "someHash", "timeStamp")
		err := db.AddUser("user@example.com", "someHash", "timeStamp")
		if err == nil {
			t.Fatalf("add duplicate: got nil; want %v", domain.ErrEmailAlreadyExists)
		}
		if !errors.Is(err, domain.ErrEmailAlreadyExists) {
			t.Errorf("add duplicate: got %v; want %v", err, domain.ErrEmailAlreadyExists)
		}
	})

	t.Run("concurrent duplicates", func(t *testing.T) {
		db := storage.NewDatabase()
		var wg sync.WaitGroup
		result := make([]error, 5)
		wg.Add(5)
		for i := 0; i < 5; i++ {
			go func(j int) {
				defer wg.Done()
				err := db.AddUser("same@example.com", "hash", "timeStamp")
				if err != nil {
					result[j] = err
				}
			}(i)
		}
		wg.Wait()
		nilCounter := 0
		errCounter := 0
		for _, v := range result {
			if v == nil {
				nilCounter++
				continue
			}
			if errors.Is(v, domain.ErrEmailAlreadyExists) {
				errCounter++
			} else {
				t.Errorf("add duplicate: got %v; want %v", v, domain.ErrEmailAlreadyExists)
			}
		}
		if nilCounter != 1 || errCounter != 4 {
			t.Errorf("thread safety: got %d-%d; want %d-%d", nilCounter, errCounter, 1, 4)
		}

	})

}
