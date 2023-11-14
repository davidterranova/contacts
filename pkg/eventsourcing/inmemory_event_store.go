package eventsourcing

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type eventStore[T Aggregate] struct {
	aggregateEvents   map[uuid.UUID][]Event[T]
	unPublishedEvents []Event[T]
}

func NewInMemoryEventStore[T Aggregate]() *eventStore[T] {
	return &eventStore[T]{
		aggregateEvents:   make(map[uuid.UUID][]Event[T]),
		unPublishedEvents: make([]Event[T], 0),
	}
}

func (s *eventStore[T]) Store(_ context.Context, events ...Event[T]) error {
	for _, event := range events {
		//nolint:staticcheck
		localEvents, ok := s.aggregateEvents[event.AggregateId()]
		if !ok {
			//nolint:staticcheck
			localEvents = make([]Event[T], 0)
		}
		localEvents = append(s.aggregateEvents[event.AggregateId()], event)
		s.unPublishedEvents = append(s.unPublishedEvents, event)

		s.aggregateEvents[event.AggregateId()] = localEvents
		log.Debug().Str("type", event.EventType()).Interface("event", event).Msg("stored event")
	}
	return nil
}

func (s *eventStore[T]) Load(_ context.Context, _ AggregateType, aggregateId uuid.UUID) ([]Event[T], error) {
	localEvents, ok := s.aggregateEvents[aggregateId]
	if !ok {
		return nil, nil
	}

	events := make([]Event[T], 0, len(localEvents))
	events = append(events, localEvents...)

	return events, nil
}

func (s *eventStore[T]) LoadUnpublished(_ context.Context, batchSize int) ([]Event[T], error) {
	events := make([]Event[T], 0, batchSize)
	added := 0
	for _, unPublished := range s.unPublishedEvents {
		log.Debug().Str("type", unPublished.EventType()).Interface("event", unPublished).Msg("loaded unpublished event")
		events = append(events, unPublished)
		added++
		if added >= batchSize {
			break
		}
	}

	return events, nil
}

func (s *eventStore[T]) MarkPublished(_ context.Context, events ...Event[T]) error {
	for _, event := range events {
		for i, unPublished := range s.unPublishedEvents {
			if unPublished.Id() == event.Id() {
				log.Debug().Str("type", event.EventType()).Interface("event", event).Msg("marked event as published")
				s.unPublishedEvents = append(s.unPublishedEvents[:i], s.unPublishedEvents[i+1:]...)
				break
			}
		}
	}

	return nil
}
