package app

import (
	"auth-service/internal/cli"
	"auth-service/internal/domain"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"errors"
	"fmt"
	"io"
	"os"
)

func Run() {
	db := storage.NewDatabase()
	emailAddress, password, err := cli.ReadRegistrationData(db)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Println("\nEOF")
			return
		}
		fmt.Fprintf(os.Stderr, "in inputFunc: %v\n", err)
		return
	}
	if err := service.Register(emailAddress, password, db); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			fmt.Fprintf(os.Stderr, "in register: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}
}
