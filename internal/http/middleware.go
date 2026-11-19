package http

import (
	"log"
	"net/http"
)

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("beginning of: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("end of: %s %s", r.Method, r.URL.Path)
	})
}
