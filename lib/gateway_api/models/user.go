package models

import "context"

const userCtxKey ctxKey = "userKey"

func UserFromCtx(ctx context.Context) *User {
	if user, ok := ctx.Value(userCtxKey).(User); ok {
		return &user
	}

	return nil
}

func AddUserToContext(ctx context.Context, user User) {
	ctx = context.WithValue(ctx, userCtxKey, user)
}

type UserService interface {
	GetUser(ctx context.Context) (*User, error)
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
