package domain

import (
	"github.com/davidterranova/contacts/pkg/user"

	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
)

type EvtContactCreated struct {
	eventsourcing.EventBase[*Contact]
}

func NewEvtContactCreated(aggregateId uuid.UUID, createdBy user.User) EvtContactCreated {
	return EvtContactCreated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, createdBy, aggregateId),
	}
}

func (e EvtContactCreated) EventType() string {
	return "contact.created"
}

func (e EvtContactCreated) Apply(contact *Contact) error {
	contact.Id = e.AggregateId()
	contact.CreatedAt = e.CreatedAt()
	contact.CreatedBy = e.IssuedBy()

	return nil
}

type EvtContactEmailUpdated struct {
	eventsourcing.EventBase[*Contact]

	email string
}

func NewEvtContactEmailUpdated(aggregateId uuid.UUID, updatedBy user.User, email string) EvtContactEmailUpdated {
	return EvtContactEmailUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, updatedBy, aggregateId),
		email:     email,
	}
}

func (e EvtContactEmailUpdated) EventType() string {
	return "contact.updated-email"
}

func (e EvtContactEmailUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.CreatedAt()
	contact.Email = e.email

	return nil
}

type EvtContactNameUpdated struct {
	eventsourcing.EventBase[*Contact]

	firstName string
	lastName  string
}

func NewEvtContactNameUpdated(aggregateId uuid.UUID, updatedBy user.User, firstName string, lastName string) EvtContactNameUpdated {
	return EvtContactNameUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, updatedBy, aggregateId),
		firstName: firstName,
		lastName:  lastName,
	}
}

func (e EvtContactNameUpdated) EventType() string {
	return "contact.updated-name"
}

func (e EvtContactNameUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.CreatedAt()
	contact.FirstName = e.firstName
	contact.LastName = e.lastName

	return nil
}

type EvtContactPhoneUpdated struct {
	eventsourcing.EventBase[*Contact]

	phone string
}

func NewEvtContactPhoneUpdated(aggregateId uuid.UUID, updatedBy user.User, phone string) EvtContactPhoneUpdated {
	return EvtContactPhoneUpdated{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, updatedBy, aggregateId),
		phone:     phone,
	}
}

func (e EvtContactPhoneUpdated) EventType() string {
	return "contact.updated-phone"
}

func (e EvtContactPhoneUpdated) Apply(contact *Contact) error {
	contact.UpdatedAt = e.CreatedAt()
	contact.Phone = e.phone

	return nil
}

type EvtContactDeleted struct {
	eventsourcing.EventBase[*Contact]
}

func NewEvtContactDeleted(aggregateId uuid.UUID, deletedBy user.User) EvtContactDeleted {
	return EvtContactDeleted{
		EventBase: eventsourcing.NewEventBase[*Contact](AggregateContact, deletedBy, aggregateId),
	}
}

func (e EvtContactDeleted) EventType() string {
	return "contact.deleted"
}

func (e EvtContactDeleted) Apply(contact *Contact) error {
	deletedAt := e.CreatedAt()
	contact.DeletedAt = &deletedAt

	return nil
}
