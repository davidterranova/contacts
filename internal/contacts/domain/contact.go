package domain

import (
	"time"

	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/cqrs/eventsourcing"
	"github.com/google/uuid"
)

const AggregateContact eventsourcing.AggregateType = "contact"

type Contact struct {
	*eventsourcing.AggregateBase[Contact]

	DeletedAt *time.Time
	CreatedBy user.User

	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func New() *Contact {
	return &Contact{
		AggregateBase: eventsourcing.NewAggregateBase[Contact](uuid.Nil, 0),
	}
}

func (c Contact) AggregateType() eventsourcing.AggregateType {
	return AggregateContact
}
