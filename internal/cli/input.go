package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kimkachailo/auth-service/internal/domain"
	"github.com/kimkachailo/auth-service/internal/email"
)

type EmailFinder interface {
	SearchEmail(string) (domain.User, error)
}

func ReadRegistrationData(db EmailFinder) (string, string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your e-mail: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", "", err
			} else {
				return "", "", io.EOF
			}
		}
		validatedEmailAddress, err := email.Validate(scanner.Text())
		if err != nil {
			if errors.Is(err, email.ErrInvalidEmail) {
				fmt.Println("Invalid e-mail, please try again")
				continue
			} else {
				return "", "", err
			}
		}

		if _, err = db.SearchEmail(validatedEmailAddress); err == nil {
			fmt.Println("This e-mail already exists")
			continue
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			return "", "", err
		}

		fmt.Print("Enter your password: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", "", err
			} else {
				return "", "", io.EOF
			}
		}
		return validatedEmailAddress, scanner.Text(), nil
	}
}
