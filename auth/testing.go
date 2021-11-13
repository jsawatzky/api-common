package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/jsawatzky/go-common/api"
)

func NewTestingMiddleware(ip IdentityProvider, uid string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			user, err := ip.Get(r.Context(), uid)
			if err != nil {
				if errors.Is(err, ErrPermissionDenied) {
					logger.Warn("Permission denied for user \"%s\"", uid)
				} else {
					logger.Error("Error retreiving user identity: %v", err)
				}
				api.EncodeResponse(rw, http.StatusForbidden, AuthError(err.Error()))
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			h.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
