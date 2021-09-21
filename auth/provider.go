package auth

import (
	"context"
	"errors"
)

var (
	ErrPermissionDenied = errors.New("user not allowed")
)

type IdentityProvider interface {
	Get(ctx context.Context, id string) (interface{}, error)
}
