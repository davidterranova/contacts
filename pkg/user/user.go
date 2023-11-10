package user

import (
	"encoding/json"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeSystem          UserType = "system"
	UserTypeAuthenticated   UserType = "authenticated"
	UserTypeUnauthenticated UserType = "unauthenticated"
)

var (
	Unauthenticated = new(uuid.Nil, UserTypeUnauthenticated)
	System          = new(uuid.Nil, UserTypeSystem)
)

type User struct {
	id       uuid.UUID
	userType UserType
}

func New(id uuid.UUID) User {
	return new(id, UserTypeAuthenticated)
}

func new(id uuid.UUID, userType UserType) User {
	return User{
		id:       id,
		userType: userType,
	}
}

func (u User) Id() uuid.UUID {
	return u.id
}

func (u User) Type() UserType {
	return u.userType
}

func (u User) IsAuthenticatedOrSystem() bool {
	return u.userType == UserTypeAuthenticated || u.userType == UserTypeSystem
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id   uuid.UUID `json:"id"`
		Type string    `json:"type"`
	}{
		Id:   u.id,
		Type: string(u.userType),
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type alias struct {
		Id   uuid.UUID `json:"id"`
		Type string    `json:"type"`
	}

	var a alias
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	u.id = a.Id
	u.userType = UserType(a.Type)

	return nil
}
