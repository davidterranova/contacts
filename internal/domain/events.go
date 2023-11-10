package domain

import (
	"github.com/davidterranova/contacts/pkg/user"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
)

const (
	ContactCreated      = "created"
	ContactEmailUpdated = "updated-email"
	ContactNameUpdated  = "updated-name"
	ContactPhoneUpdated = "updated-phone"
	ContactDeleted      = "deleted"
)

type EvtContactCreated struct {
	eventsourcing.EventBase[*Contact]
}

func NewEvtContactCreated(aggregateId uuid.UUID, createdBy user.User) *EvtContactCreated {
	return &EvtContactCreated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, ContactCreated, aggregateId, createdBy),
	}
}

func (e EvtContactCreated) Apply(contact *Contact) error {
	contact.Id = e.AggregateId()
	contact.CreatedAt = e.IssuedAt()
	contact.UpdatedAt = e.IssuedAt()
	contact.CreatedBy = e.IssuedBy()

	return nil
}

type EvtContactEmailUpdated struct {
	eventsourcing.EventBase[*Contact]

	email string
}

func NewEvtContactEmailUpdated(aggregateId uuid.UUID, updatedBy user.User, email string) *EvtContactEmailUpdated {
	return &EvtContactEmailUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, ContactEmailUpdated, aggregateId, updatedBy),
		email:     email,
	}
}

func (e EvtContactEmailUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.IssuedAt()
	contact.Email = e.email

	return nil
}

type EvtContactNameUpdated struct {
	eventsourcing.EventBase[*Contact]

	firstName string
	lastName  string
}

func NewEvtContactNameUpdated(aggregateId uuid.UUID, updatedBy user.User, firstName string, lastName string) *EvtContactNameUpdated {
	return &EvtContactNameUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, ContactNameUpdated, aggregateId, updatedBy),
		firstName: firstName,
		lastName:  lastName,
	}
}

func (e EvtContactNameUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.IssuedAt()
	contact.FirstName = e.firstName
	contact.LastName = e.lastName

	return nil
}

type EvtContactPhoneUpdated struct {
	eventsourcing.EventBase[*Contact]

	phone string
}

func NewEvtContactPhoneUpdated(aggregateId uuid.UUID, updatedBy user.User, phone string) *EvtContactPhoneUpdated {
	return &EvtContactPhoneUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, ContactPhoneUpdated, aggregateId, updatedBy),
		phone:     phone,
	}
}

func (e EvtContactPhoneUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.IssuedAt()
	contact.Phone = e.phone

	return nil
}

type EvtContactDeleted struct {
	eventsourcing.EventBase[*Contact]
}

func NewEvtContactDeleted(aggregateId uuid.UUID, deletedBy user.User) *EvtContactDeleted {
	return &EvtContactDeleted{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, ContactDeleted, aggregateId, deletedBy),
	}
}

func (e EvtContactDeleted) Apply(contact *Contact) error {
	deletedAt := e.IssuedAt()
	contact.DeletedAt = &deletedAt

	return nil
}
