package auth

import (
	"context"
	"errors"

	"github.com/davidterranova/contacts/internal/domain"
)

type RequestCtxKey string

const RequestCtxUserKey RequestCtxKey = "user"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("unauthorized")
)

func UserFromContext(ctx context.Context) (domain.User, error) {
	u, ok := ctx.Value(RequestCtxUserKey).(domain.User)
	if !ok {
		return *domain.NewEmptyUser(), ErrUserNotFound
	}

	return u, nil
}

func ContextWithUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, RequestCtxUserKey, u)
}
