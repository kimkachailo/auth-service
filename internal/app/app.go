package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kimkachailo/auth-service/internal/cli"
	"github.com/kimkachailo/auth-service/internal/domain"
	myhttp "github.com/kimkachailo/auth-service/internal/http"
	"github.com/kimkachailo/auth-service/internal/service"
	"github.com/kimkachailo/auth-service/internal/storage"
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
