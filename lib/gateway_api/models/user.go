package models

const UserCtxKey = "user-ctx-key"

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
