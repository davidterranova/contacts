package user

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	pkguser "github.com/davidterranova/cqrs/user"
	"github.com/google/uuid"
)

var (
	Unauthenticated = new(uuid.Nil, pkguser.UserTypeUnauthenticated)
	System          = new(uuid.Nil, pkguser.UserTypeSystem)
)

type User interface {
	Id() uuid.UUID
	Type() pkguser.UserType
	IsAuthenticatedOrSystem() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	String() string
	FromString(string) error
}

type user struct {
	id       uuid.UUID
	userType pkguser.UserType
}

func New(id uuid.UUID) *user {
	return new(id, pkguser.UserTypeAuthenticated)
}

func new(id uuid.UUID, userType pkguser.UserType) *user {
	return &user{
		id:       id,
		userType: userType,
	}
}

func (u user) Id() uuid.UUID {
	return u.id
}

func (u user) Type() pkguser.UserType {
	return u.userType
}

func (u user) IsAuthenticatedOrSystem() bool {
	return u.userType == pkguser.UserTypeAuthenticated || u.userType == pkguser.UserTypeSystem
}

func (u user) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id   uuid.UUID `json:"id"`
		Type string    `json:"type"`
	}{
		Id:   u.id,
		Type: string(u.userType),
	})
}

func (u *user) UnmarshalJSON(data []byte) error {
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
	u.userType = pkguser.UserType(a.Type)

	return nil
}

func (u user) String() string {
	return fmt.Sprintf("%s:%s", u.userType, u.id)
}

func (u *user) FromString(s string) error {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ':'
	records, err := r.Read()
	if err != nil {
		return fmt.Errorf("failed to read user string: %w", err)
	}

	if len(records) != 2 {
		return fmt.Errorf("invalid user string: %s", s)
	}

	u.userType = pkguser.UserType(records[0])
	u.id, err = uuid.Parse(records[1])

	return err
}
