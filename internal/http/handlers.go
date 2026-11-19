package http

import (
	"auth-service/internal/email"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
)

func health(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("write health response: %v", err)

	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	mediaT, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaT != "application/json" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	var input struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		http.Error(w, "expected one JSON value", http.StatusBadRequest)
		return
	}

	validatedEmail, err := email.Validate(input.Email)
	if err != nil {
		http.Error(w, email.ErrInvalidEmail.Error(), http.StatusBadRequest)
		return
	}
	input.Email = validatedEmail

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(input); err != nil {
		log.Printf("write echo response: %v", err)
	}

}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go") // add header field
	w.Write([]byte("Welcome to auth-service"))
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil || userID < 0 {
		http.NotFound(w, r) // response with 404 error
		return
	}
	fmt.Fprintf(w, "Information about user with ID: %d", userID) // http.ResponseWriter implemets io.Writer. w can be passed to io.Writer
}

func postSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // changed header field, json type should be set explicitly as it's defined as string automaticly
	w.Write([]byte(`{"email" : "test@gmail.com"}`))    // `...` is for raw string literal, not need to shield and use \n
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil || userID < 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "User with ID: %d was deleted", userID)
}
