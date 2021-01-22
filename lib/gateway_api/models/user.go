package models

import (
	"context"
	"fmt"
)

const UserCtxKey = "user-ctx-key"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserFromContext(ctx context.Context) (User, error) {
	if user, ok := ctx.Value(UserCtxKey).(User); ok {
		return user, nil
	}

	return User{}, fmt.Errorf("user not found in context")
}
