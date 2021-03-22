package models

import (
	"context"
	"fmt"
)

const userCtxKey ctxKey = "userKey"

func GetUserFromContext(ctx context.Context) (*User, error) {
	if user, ok := ctx.Value(userCtxKey).(User); ok {
		return &user, nil
	}

	return nil, fmt.Errorf("user not found in context")
}

func AddUserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

type UserService interface {
	GetUser(ctx context.Context) (*User, error)
}

type User struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}
