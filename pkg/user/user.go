package user

import "github.com/google/uuid"

type UserType string

const (
	UserTypeSystem        UserType = "system"
	UserTypeAuthenticated UserType = "authenticated"
	UserTypeUnknown       UserType = "unknown"
)

type User interface {
	Id() uuid.UUID
	Type() UserType
}

func New(id uuid.UUID, userType UserType) User {
	switch userType {
	case UserTypeAuthenticated:
		return &UserAuthenticated{id: id}
	case UserTypeSystem:
		return &UserSystem{id: id}
	default:
		return &UserUnknown{}
	}
}

type UserAuthenticated struct {
	id uuid.UUID
}

func (u UserAuthenticated) Id() uuid.UUID {
	return u.id
}

func (u UserAuthenticated) Type() UserType {
	return UserTypeAuthenticated
}

type UserSystem struct {
	id uuid.UUID
}

func (u UserSystem) Id() uuid.UUID {
	return u.id
}

func (u UserSystem) Type() UserType {
	return UserTypeSystem
}

type UserUnknown struct{}

func (u UserUnknown) Id() uuid.UUID {
	return uuid.Nil
}

func (u UserUnknown) Type() UserType {
	return UserTypeUnknown
}

func NewEmpty() User {
	return &UserUnknown{}
}
