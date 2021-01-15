package middleware

import (
	"log"
	"net/http"
)

// RequestLogger Middleware logs requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.RequestURI, r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

// ContentType Middleware which sets default content type as application/json
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
