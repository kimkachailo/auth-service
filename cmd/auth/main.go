package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrEmailAlreadyExists = errors.New("Email already exists")
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
	if _, err := d.SearchEmail(email); err == nil {
		return ErrEmailAlreadyExists
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

func hashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return string(h.Sum(nil))
}

func checkEmailForCorrectness(email string) bool {
	domains := []string{"@gmail.com", "@outlook.com", "@hotmail.com", "@live.com", "@yahoo.com", "@icloud.com", "@proton.me", "@protonmail.com", "@zoho.com", "@aol.com", "@mail.ru", "@inbox.ru", "@list.ru", "@bk.ru", "@yandex.ru", "@ya.ru", "@rambler.ru", "@lenta.ru"}
	var flag bool
	for i := range domains {
		if strings.HasSuffix(email, domains[i]) {
			flag = true
		}
	}
	return flag
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

		email = scanner.Text()
		if !checkEmailForCorrectness(email) {
			fmt.Println("Incorrect e-mail, please try again")
			continue
		}

		if _, err := db.SearchEmail(email); err == nil {
			fmt.Println(ErrEmailAlreadyExists)
			continue
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

		passHash = hashPassword(scanner.Text())
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
