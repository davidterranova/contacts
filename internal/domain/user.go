package domain

import "github.com/google/uuid"

type User struct {
	Id uuid.UUID
}

func NewUser(id uuid.UUID) *User {
	return &User{Id: id}
}

func NewEmptyUser() *User {
	return &User{}
}
