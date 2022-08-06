package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/jsawatzky/go-common/api"
)

func NewReverseProxyMiddleware(ip IdentityProvider) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("Remote-User")
			user, err := ip.Get(r.Context(), id)
			if err != nil {
				if errors.Is(err, ErrPermissionDenied) {
					logger.Warn("Permission denied for user \"%s\"", id)
				} else {
					logger.Error("Error retrieving user identity: %v", err)
				}
				api.EncodeResponse(rw, http.StatusForbidden, AuthError(err.Error()))
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			h.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
