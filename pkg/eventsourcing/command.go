package eventsourcing

import (
	"time"

	"github.com/google/uuid"
)

type Command[T Aggregate] interface {
	AggregateId() uuid.UUID
	AggregateType() AggregateType
	CreatedAt() time.Time

	// Check for validity of command on aggregate, mutate the aggregate and return newly emitted events
	Apply(T) ([]Event[T], error)
}

type BaseCommand[T Aggregate] struct {
	aggregateId   uuid.UUID     `validate:"required"`
	aggregateType AggregateType `validate:"required"`
	createdAt     time.Time     `validate:"required"`
}

func NewBaseCommand[T Aggregate](aggregateId uuid.UUID, aggregateType AggregateType) BaseCommand[T] {
	return BaseCommand[T]{
		aggregateId:   aggregateId,
		aggregateType: aggregateType,
		createdAt:     time.Now().UTC(),
	}
}

func (c BaseCommand[T]) AggregateId() uuid.UUID {
	return c.aggregateId
}

func (c BaseCommand[T]) AggregateType() AggregateType {
	return c.aggregateType
}

func (c BaseCommand[T]) CreatedAt() time.Time {
	return c.createdAt
}
