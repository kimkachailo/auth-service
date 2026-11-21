package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kimkachailo/auth-service/internal/domain"
	"github.com/kimkachailo/auth-service/internal/email"
	"github.com/kimkachailo/auth-service/internal/hash"
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

		fmt.Println("Password hint. Use 8–72 characters (Latin letters, digits, and !@#$%^&*). Include at least one uppercase letter, one digit, and one symbol.")
		fmt.Print("Enter your password: ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", "", err
			} else {
				return "", "", io.EOF
			}
		}
		validatedPassword, err := hash.ValidatePassword(scanner.Text())
		if err != nil {
			if errors.Is(err, hash.ErrInvalidPassword) {
				fmt.Println("Password does not meet the requirements")
				continue
			}
			return "", "", err
		}
		return validatedEmailAddress, validatedPassword, nil
	}
}
