package domain

import (
	"time"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
)

const AggregateContact eventsourcing.AggregateType = "contact"

type Contact struct {
	*eventsourcing.AggregateBase[Contact]

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
		AggregateBase: eventsourcing.NewAggregateBase[Contact](uuid.Nil),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (c Contact) AggregateType() eventsourcing.AggregateType {
	return AggregateContact
}
