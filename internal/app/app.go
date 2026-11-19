package app

import (
	"auth-service/internal/cli"
	"auth-service/internal/domain"
	myhttp "auth-service/internal/http"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func RunCLI() {
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

func RunServer() {
	handler := myhttp.LogRequests(myhttp.NewRouter())
	log.Print("starting server on :8080")
	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}
