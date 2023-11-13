package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

type EventStore[T Aggregate] interface {
	Store(ctx context.Context, events ...Event[T]) error
	Load(aggregateType AggregateType, aggregateId uuid.UUID) ([]Event[T], error)
}

type eventStore[T Aggregate] struct {
	storage map[uuid.UUID][]Event[T]
}

func NewInMemoryEventStore[T Aggregate]() *eventStore[T] {
	return &eventStore[T]{
		storage: make(map[uuid.UUID][]Event[T]),
	}
}

func (s *eventStore[T]) Store(events ...Event[T]) error {
	for _, event := range events {
		//nolint:staticcheck
		localEvents, ok := s.storage[event.AggregateId()]
		if !ok {
			//nolint:staticcheck
			localEvents = make([]Event[T], 0)
		}
		localEvents = append(s.storage[event.AggregateId()], event)

		s.storage[event.AggregateId()] = localEvents
	}
	return nil
}

func (s *eventStore[T]) Load(_ AggregateType, aggregateId uuid.UUID) ([]Event[T], error) {
	events, ok := s.storage[aggregateId]
	if !ok {
		return nil, nil
	}

	return events, nil
}
