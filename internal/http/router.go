package http

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", welcomePage)
	mux.HandleFunc("GET /users/{userID}", getUserInfo)
	mux.HandleFunc("POST /signup", postSignUp)
	mux.HandleFunc("DELETE /users/{userID}", deleteUser)
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /echo", echo)
	return mux
}
