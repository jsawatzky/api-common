package recovery

import (
	"net/http"

	"github.com/jsawatzky/go-common/log"
)

var (
	logger = log.GetLogger("recovery")
)

func recoverFromPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		logger.Error("unhandled panic: %v", r)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer recoverFromPanic(rw)
		h.ServeHTTP(rw, r)
	})
}
