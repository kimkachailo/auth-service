package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    string
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

func (d *Database) AddUser(email string, hash string, timeStamp string) {
	d.CounterID++
	d.Storage[d.CounterID] = User{
		ID:           d.CounterID,
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    timeStamp,
	}
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

func inputFunc() (email, passHash, timeStamp string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your e-mail: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			return
		}
		email = scanner.Text()
		if !checkEmailForCorrectness(email) {
			fmt.Println("Incorrect e-mail, please try again")
			continue
		}

		fmt.Print("Enter your password: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			return
		}
		passHash = hashPassword(scanner.Text())
		break
	}
	timeStamp = time.Now().UTC().Truncate(time.Second).String()
	return
}

func main() {
	db := NewDatabase()
	db.AddUser(inputFunc())
}
