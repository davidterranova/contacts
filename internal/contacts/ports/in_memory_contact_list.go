package ports

import (
	"context"
	"errors"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/cqrs/eventsourcing"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var ErrUnknownEvent = errors.New("unknown event")

type InMemoryContactList struct {
	contacts map[uuid.UUID]*domain.Contact
}

func NewInMemoryContactList(eventStream eventsourcing.Subscriber[domain.Contact]) *InMemoryContactList {
	l := &InMemoryContactList{
		contacts: map[uuid.UUID]*domain.Contact{},
	}
	eventStream.Subscribe(context.Background(), l.HandleEvent)

	return l
}

func (l *InMemoryContactList) HandleEvent(e eventsourcing.Event[domain.Contact]) {
	var (
		err error
		c   *domain.Contact
		ok  bool
	)

	switch e.EventType() {
	case domain.ContactCreated:
		c = domain.New()
		err = e.Apply(c)
	case domain.ContactEmailUpdated, domain.ContactNameUpdated, domain.ContactPhoneUpdated:
		c, ok = l.contacts[e.AggregateId()]
		if !ok {
			log.Error().Err(ErrUnknownEvent).Msgf("event %s for unknown contact %s", e.EventType(), e.AggregateId())
			return
		}
		err = e.Apply(c)
	case domain.ContactDeleted:
		delete(l.contacts, e.AggregateId())
		return
	default:
		log.Error().Err(ErrUnknownEvent).Msgf("unknown event %s", e.EventType())
		return
	}

	if err != nil {
		log.Error().Err(err).Msgf("error applying event %s on contact %q", e.EventType(), e.AggregateId())
	}

	l.contacts[e.AggregateId()] = c
}

func (l *InMemoryContactList) List(_ context.Context, query usecase.QueryContact) ([]*domain.Contact, error) {
	contacts := make([]*domain.Contact, 0, len(l.contacts))
	for _, contact := range l.contacts {
		if l.matchQuery(contact, query) {
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}

func (l *InMemoryContactList) Get(_ context.Context, query usecase.QueryContact) (*domain.Contact, error) {
	for _, contact := range l.contacts {
		if l.matchQuery(contact, query) {
			return contact, nil
		}
	}

	return nil, ErrNotFound
}

func (l *InMemoryContactList) matchQuery(contact *domain.Contact, query usecase.QueryContact) bool {
	toAdd := true
	if contact.CreatedBy.Id() != query.Requestor.Id() {
		toAdd = false
	}

	if query.ContactId != nil && contact.AggregateId() != *query.ContactId {
		toAdd = false
	}

	return toAdd
}
