package domain

import (
	luser "github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/cqrs/user"

	"github.com/davidterranova/cqrs/eventsourcing"
	"github.com/google/uuid"
)

const (
	ContactCreated      = "created"
	ContactEmailUpdated = "updated-email"
	ContactNameUpdated  = "updated-name"
	ContactPhoneUpdated = "updated-phone"
	ContactDeleted      = "deleted"
)

func RegisterEvents(registry eventsourcing.EventRegistry[Contact]) {
	registry.Register(ContactCreated, func() eventsourcing.Event[Contact] {
		return &EvtContactCreated{
			EventBase: &eventsourcing.EventBase[Contact]{},
		}
	})
	registry.Register(ContactEmailUpdated, func() eventsourcing.Event[Contact] {
		return &EvtContactEmailUpdated{
			EventBase: &eventsourcing.EventBase[Contact]{},
		}
	})
	registry.Register(ContactNameUpdated, func() eventsourcing.Event[Contact] {
		return &EvtContactNameUpdated{
			EventBase: &eventsourcing.EventBase[Contact]{},
		}
	})
	registry.Register(ContactPhoneUpdated, func() eventsourcing.Event[Contact] {
		return &EvtContactPhoneUpdated{
			EventBase: &eventsourcing.EventBase[Contact]{},
		}
	})
	registry.Register(ContactDeleted, func() eventsourcing.Event[Contact] {
		return &EvtContactDeleted{
			EventBase: &eventsourcing.EventBase[Contact]{},
		}
	})
}

type EvtContactCreated struct {
	*eventsourcing.EventBase[Contact]
}

func NewEvtContactCreated(aggregateId uuid.UUID, aggregateVersion int, createdBy user.User) *EvtContactCreated {
	return &EvtContactCreated{
		EventBase: eventsourcing.NewEventBase[Contact](
			AggregateContact,
			aggregateVersion,
			ContactCreated,
			aggregateId,
			createdBy,
		),
	}
}

func (e EvtContactCreated) Apply(contact *Contact) error {
	contact.Init(e)
	lu := e.IssuedBy().(luser.User)
	contact.CreatedBy = lu

	return nil
}

type EvtContactEmailUpdated struct {
	*eventsourcing.EventBase[Contact]

	Email string `json:"email"`
}

func NewEvtContactEmailUpdated(aggregateId uuid.UUID, aggregateVersion int, updatedBy user.User, email string) *EvtContactEmailUpdated {
	return &EvtContactEmailUpdated{
		EventBase: eventsourcing.NewEventBase[Contact](
			AggregateContact,
			aggregateVersion,
			ContactEmailUpdated,
			aggregateId,
			updatedBy,
		),
		Email: email,
	}
}

func (e EvtContactEmailUpdated) Apply(contact *Contact) error {
	contact.Process(e)

	contact.UpdatedAt = e.IssuedAt()
	contact.Email = e.Email

	return nil
}

type EvtContactNameUpdated struct {
	*eventsourcing.EventBase[Contact]

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewEvtContactNameUpdated(aggregateId uuid.UUID, aggregateVersion int, updatedBy user.User, firstName string, lastName string) *EvtContactNameUpdated {
	return &EvtContactNameUpdated{
		EventBase: eventsourcing.NewEventBase[Contact](
			AggregateContact,
			aggregateVersion,
			ContactNameUpdated,
			aggregateId,
			updatedBy,
		),
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (e EvtContactNameUpdated) Apply(contact *Contact) error {
	contact.Process(e)

	contact.UpdatedAt = e.IssuedAt()
	contact.FirstName = e.FirstName
	contact.LastName = e.LastName

	return nil
}

type EvtContactPhoneUpdated struct {
	*eventsourcing.EventBase[Contact]

	Phone string `json:"phone"`
}

func NewEvtContactPhoneUpdated(aggregateId uuid.UUID, aggregateVersion int, updatedBy user.User, phone string) *EvtContactPhoneUpdated {
	return &EvtContactPhoneUpdated{
		EventBase: eventsourcing.NewEventBase[Contact](
			AggregateContact,
			aggregateVersion,
			ContactPhoneUpdated,
			aggregateId,
			updatedBy,
		),
		Phone: phone,
	}
}

func (e EvtContactPhoneUpdated) Apply(contact *Contact) error {
	contact.Process(e)

	contact.UpdatedAt = e.IssuedAt()
	contact.Phone = e.Phone

	return nil
}

type EvtContactDeleted struct {
	*eventsourcing.EventBase[Contact]
}

func NewEvtContactDeleted(aggregateId uuid.UUID, aggregateVersion int, deletedBy user.User) *EvtContactDeleted {
	return &EvtContactDeleted{
		EventBase: eventsourcing.NewEventBase[Contact](
			AggregateContact,
			aggregateVersion,
			ContactDeleted,
			aggregateId,
			deletedBy,
		),
	}
}

func (e EvtContactDeleted) Apply(contact *Contact) error {
	contact.Process(e)

	deletedAt := e.IssuedAt()
	contact.DeletedAt = &deletedAt

	return nil
}
