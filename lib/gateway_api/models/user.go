package models

import "context"

type UserService interface {
	GetUser(ctx context.Context) (*User, error)
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
