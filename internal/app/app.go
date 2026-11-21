package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kimkachailo/auth-service/internal/cli"
	myhttp "github.com/kimkachailo/auth-service/internal/http"
	"github.com/kimkachailo/auth-service/internal/service"
	"github.com/kimkachailo/auth-service/internal/storage"
)

func RunCLI() {
	db := storage.NewDatabase()
	emailAddress, password, err := cli.ReadRegistrationData(db)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr, "\nInput ended before registration was complete")
			return
		}
		fmt.Fprintf(os.Stderr, "in inputFunc: %v\n", err)
		return
	}
	if err := service.SignUp(emailAddress, password, db); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func RunServer() {
	handler := myhttp.LogRequests(myhttp.NewRouter())
	log.Print("starting server on :8080")
	err := http.ListenAndServe(":8080", handler)
	log.Fatal(err)
}
