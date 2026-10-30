package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"time"
)

var Database = make(map[int]User)
var CounterID int

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    string
}

func CreateUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, passHash string
	domains := []string{"@gmail.com", "@outlook.com", "@hotmail.com", "@live.com", "@yahoo.com", "@icloud.com", "@proton.me", "@protonmail.com", "@zoho.com", "@aol.com", "@mail.ru", "@inbox.ru", "@list.ru", "@bk.ru", "@yandex.ru", "@ya.ru", "@rambler.ru", "@lenta.ru"}
	for {
		fmt.Print("Enter your e-mail: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			return
		}
		email = scanner.Text()
		var flag bool
		for i := range domains {
			if strings.HasSuffix(email, domains[i]) {
				flag = true
			}
		}
		if !flag {
			fmt.Println("Incorrect e-mail, please try again")
			continue
		}

		fmt.Print("Enter your password: ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			return
		}
		h := sha256.New()
		h.Write([]byte(scanner.Text()))
		passHash = string(h.Sum(nil))
		break
	}
	timeStamp := time.Now().UTC().Truncate(time.Second).String()
	CounterID++

	Database[CounterID] = User{
		ID:           CounterID,
		Email:        email,
		PasswordHash: passHash,
		CreatedAt:    timeStamp,
	}
	fmt.Println(Database)
}

func main() {
	CreateUser()
}
