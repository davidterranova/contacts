package ports

import (
	"context"
	"errors"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrUnknownEvent = errors.New("unknown event")

	ContactCreated      = domain.EvtContactCreated{}.EventType()
	ContactEmailUpdated = domain.EvtContactEmailUpdated{}.EventType()
	ContactNameUpdated  = domain.EvtContactNameUpdated{}.EventType()
	ContactPhoneUpdated = domain.EvtContactPhoneUpdated{}.EventType()
	ContactDeleted      = domain.EvtContactDeleted{}.EventType()
)

type InMemoryContactList struct {
	contacts map[uuid.UUID]*domain.Contact
}

func NewInMemoryContactList(eventStream eventsourcing.Subscriber[*domain.Contact]) *InMemoryContactList {
	l := &InMemoryContactList{
		contacts: map[uuid.UUID]*domain.Contact{},
	}
	eventStream.Subscribe(l.HandleEvent)

	return l
}

func (l *InMemoryContactList) HandleEvent(e eventsourcing.Event[*domain.Contact]) {
	var (
		err error
		c   *domain.Contact
		ok  bool
	)

	switch e.EventType() {
	case ContactCreated:
		c = domain.New()
		err = e.Apply(c)
	case ContactEmailUpdated, ContactNameUpdated, ContactPhoneUpdated:
		c, ok = l.contacts[e.AggregateId()]
		if !ok {
			log.Error().Err(ErrUnknownEvent).Msgf("event %s for unknown contact %s", e.EventType(), e.AggregateId())
			return
		}
		err = e.Apply(c)
	case ContactDeleted:
		delete(l.contacts, e.AggregateId())
		return
	default:
		log.Error().Err(ErrUnknownEvent).Msgf("unknown event %s", e.EventType())
	}

	if err != nil {
		log.Error().Err(err).Msgf("error applying event %s on contact %q", e.EventType(), e.AggregateId())
	}

	l.contacts[e.AggregateId()] = c
}

func (l *InMemoryContactList) List(_ context.Context) ([]*domain.Contact, error) {
	contacts := make([]*domain.Contact, 0, len(l.contacts))
	for _, contact := range l.contacts {
		contacts = append(contacts, contact)
	}

	return contacts, nil
}
