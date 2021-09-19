package logging

import (
	"log"
	"net/http"
	"time"

	"github.com/jsawatzky/go-common/internal"
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		resp := internal.RecordResponse(rw)
		h.ServeHTTP(resp, r)

		log.Printf("http request: \"%s %s %s\" %d %v %d", r.Method, r.URL.Path, r.Proto, resp.Status(), time.Since(start), resp.ResponseSize())
	})
}
