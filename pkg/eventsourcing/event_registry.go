package eventsourcing

import (
	"encoding/json"
	"fmt"
)

type Registry[T Aggregate] struct {
	registry map[string]func() Event[T]
}

func NewRegistry[T Aggregate]() *Registry[T] {
	return &Registry[T]{
		registry: make(map[string]func() Event[T]),
	}
}

func (r *Registry[T]) Register(eventType string, factory func() Event[T]) {
	r.registry[eventType] = factory
}

func (r Registry[T]) create(eventType string) (Event[T], error) {
	factory, ok := r.registry[eventType]
	if !ok {
		return nil, ErrUnknownEventType
	}

	return factory(), nil
}

func (r Registry[T]) Hydrate(base EventBase[T], data []byte) (Event[T], error) {
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
