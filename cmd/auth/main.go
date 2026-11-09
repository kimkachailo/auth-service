package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("invalid email")
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    string
}

type EmailFinder interface {
	SearchEmail(string) (User, error)
}
type Database struct {
	Storage   map[int]User
	CounterID int
}

func NewDatabase() *Database {
	return &Database{
		Storage: make(map[int]User),
	}
}

func (d *Database) AddUser(email string, hash string, timeStamp string) error {
	_, err := d.SearchEmail(email)
	if err == nil {
		return fmt.Errorf("in AddUser: %w", ErrEmailAlreadyExists)
	}
	if !errors.Is(err, ErrUserNotFound) {
		return err
	}
	d.CounterID++
	d.Storage[d.CounterID] = User{
		ID:           d.CounterID,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    timeStamp,
	}
	return nil
}

func (d Database) SearchID(ID int) (User, error) {
	v, ok := d.Storage[ID]
	if ok {
		return v, nil
	}
	return User{}, ErrUserNotFound
}

func (d Database) SearchEmail(email string) (User, error) {
	for k, v := range d.Storage {
		if v.Email == email {
			return d.Storage[k], nil
		}
	}
	return User{}, ErrUserNotFound
}

func hashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkEmailForCorrectness(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidEmail, err)
	}
	if addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil

}

func inputFunc(db EmailFinder) (email, passHash, timeStamp string, flag bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your e-mail: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			} else {
				fmt.Println("EOF")
				return
			}
		}

		validatedEmail, err := checkEmailForCorrectness(scanner.Text())
		if err != nil {
			if errors.Is(err, ErrInvalidEmail) {
				fmt.Println("Invalid e-mail, please try again")
				continue
			} else {
				fmt.Fprintln(os.Stderr, "checkEmailForCorrectnes error: ", err)
				return
			}
		}
		email = validatedEmail

		_, err = db.SearchEmail(email)
		if err == nil {
			fmt.Println("This e-mail already exists")
			continue
		}
		if !errors.Is(err, ErrUserNotFound) {
			fmt.Fprintln(os.Stderr, "could not check email", err)
			return
		}

		fmt.Print("Enter your password: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
				return
			} else {
				fmt.Println("EOF")
				return
			}
		}

		hash, err := hashPassword(scanner.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		passHash = hash
		break
	}
	timeStamp = time.Now().UTC().Truncate(time.Second).String()
	flag = true
	return
}

func main() {
	db := NewDatabase()
	db.AddUser("test@gmail.com", "", "")
	email, passHash, timeStamp, ok := inputFunc(db)
	if !ok {
		return
	}
	if err := db.AddUser(email, passHash, timeStamp); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	userEmail, err := db.SearchEmail("test@gmail.com")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	userID, err := db.SearchID(10)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(userEmail)
	fmt.Println(userID)
}
