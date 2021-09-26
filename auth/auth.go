package auth

import "context"

type key struct{}

var (
	userKey = key{}
)

func GetUser(ctx context.Context) string {
	return ctx.Value(userKey).(string)
}
