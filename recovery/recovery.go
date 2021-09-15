package recovery

import (
	"log"
	"net/http"
)

func recoverFromPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("unhandled panic: %v", r)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer recoverFromPanic(rw)
		h.ServeHTTP(rw, r)
	})
}
