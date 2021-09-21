package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/jsawatzky/go-common/log"
)

var (
	logger = log.GetLogger("auth")
)

func NewReverseProxyMiddleware(ip IdentityProvider) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("Remote-User")
			user, err := ip.Get(r.Context(), id)
			if err != nil {
				if errors.Is(err, ErrPermissionDenied) {
					logger.Warn("Permission denied for user \"%s\"", id)
					rw.WriteHeader(http.StatusForbidden)
					return
				} else {
					logger.Error("Error retreiving user identity: %v", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			h.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
