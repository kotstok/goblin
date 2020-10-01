package errors

import (
	"net/http"
)

func NotFoundHandler() http.Handler {
	// .. 404 Http callback
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})
}

func MethodNotAllowedHandler() http.Handler {
	// .. 405 HTTP callback
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}
