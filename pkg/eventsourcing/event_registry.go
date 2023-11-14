package eventsourcing

import (
	"encoding/json"
	"fmt"
)

type EventRegistry[T Aggregate] struct {
	registry map[string]func() Event[T]
}

func NewRegistry[T Aggregate]() *EventRegistry[T] {
	return &EventRegistry[T]{
		registry: make(map[string]func() Event[T]),
	}
}

func (r *EventRegistry[T]) Register(eventType string, factory func() Event[T]) {
	r.registry[eventType] = factory
}

func (r EventRegistry[T]) create(eventType string) (Event[T], error) {
	factory, ok := r.registry[eventType]
	if !ok {
		return nil, ErrUnknownEventType
	}

	return factory(), nil
}

func (r EventRegistry[T]) Hydrate(base EventBase[T], data []byte) (Event[T], error) {
	event, err := r.create(base.EventType())
	if err != nil {
		return nil, fmt.Errorf("failed to create empty event: %w", err)
	}

	err = json.Unmarshal(data, event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}
	event.SetBase(base)

	return event, nil
}
