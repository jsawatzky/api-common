package auth

import (
	"context"

	"github.com/jsawatzky/go-common/api"
	"github.com/jsawatzky/go-common/log"
)

var (
	logger = log.GetLogger("auth")
)

type key struct{}

var (
	userKey = key{}
)

func GetUser(ctx context.Context) string {
	return ctx.Value(userKey).(string)
}

func AuthError(err string) api.Error {
	return api.Error{
		Title:   "Authorization Error",
		Details: err,
	}
}
