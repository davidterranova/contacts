package domain

import (
	"time"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
)

const AggregateContact eventsourcing.AggregateType = "contact"

type Contact struct {
	Id uuid.UUID
	*eventsourcing.AggregateBase

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy user.User

	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func New() *Contact {
	now := time.Now().UTC()

	return &Contact{
		Id:            uuid.New(),
		AggregateBase: &eventsourcing.AggregateBase{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (c Contact) AggregateId() uuid.UUID {
	return c.Id
}

func (c Contact) AggregateType() eventsourcing.AggregateType {
	return AggregateContact
}
