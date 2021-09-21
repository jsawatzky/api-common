package auth

import "context"

type key struct{}

var (
	userKey = key{}
)

func GetUser(ctx context.Context) interface{} {
	return ctx.Value(userKey)
}
