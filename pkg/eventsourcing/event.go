package eventsourcing

import (
	"time"

	"github.com/davidterranova/contacts/pkg/user"

	"github.com/google/uuid"
)

type Event[T Aggregate] interface {
	Id() uuid.UUID
	AggregateId() uuid.UUID
	AggregateType() AggregateType
	EventType() string
	CreatedAt() time.Time
	IssuedBy() user.User
	Apply(T) error
}

type EventBase[T Aggregate] struct {
	id            uuid.UUID
	aggregateType AggregateType
	aggregateId   uuid.UUID
	createdAt     time.Time
	issuedBy      user.User
}

func NewEventBase[T Aggregate](aggregateType AggregateType, issuedBy user.User, aggregateId uuid.UUID) EventBase[T] {
	return EventBase[T]{
		id:            uuid.New(),
		aggregateType: aggregateType,
		aggregateId:   aggregateId,
		issuedBy:      issuedBy,
		createdAt:     time.Now().UTC(),
	}
}

func (e EventBase[T]) Id() uuid.UUID {
	return e.id
}

func (e EventBase[T]) AggregateId() uuid.UUID {
	return e.aggregateId
}

func (e EventBase[T]) CreatedAt() time.Time {
	return e.createdAt
}

func (e EventBase[T]) AggregateType() AggregateType {
	return e.aggregateType
}

func (e EventBase[T]) IssuedBy() user.User {
	return e.issuedBy
}
