package domain

import (
	"time"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
)

const AggregateContact eventsourcing.AggregateType = "contact"

type Contact struct {
	Id uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func New() *Contact {
	now := time.Now().UTC()

	return &Contact{
		Id:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c Contact) AggregateId() uuid.UUID {
	return c.Id
}

func (c Contact) AggregateType() eventsourcing.AggregateType {
	return AggregateContact
}
