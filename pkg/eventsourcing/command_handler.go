package eventsourcing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CommandHandler[T Aggregate] interface {
	// Handle is the global command handler that should be called by the application
	Handle(ctx context.Context, cmd Command[T]) (*T, error)

	// HydrateAggregate an aggregate from already published events (internal)
	HydrateAggregate(ctx context.Context, aggregateType AggregateType, aggregateId uuid.UUID) (*T, error)

	// Apply checks command validity for an aggregate and return newly emitted events (internal)
	ApplyCommand(ctx context.Context, aggregate *T, command Command[T]) (*T, []Event[T], error)
}

type AggregateFactory[T Aggregate] func() *T

type commandHandler[T Aggregate] struct {
	eventStore     EventStore[T]
	factory        AggregateFactory[T]
	eventPublisher Publisher[T]
}

func NewCommandHandler[T Aggregate](eventStore EventStore[T], eventPublisher Publisher[T], factory AggregateFactory[T]) *commandHandler[T] {
	return &commandHandler[T]{
		eventStore:     eventStore,
		eventPublisher: eventPublisher,
		factory:        factory,
	}
}

func (h *commandHandler[T]) Handle(ctx context.Context, c Command[T]) (*T, error) {
	// hydrate aggregate
	aggregate, err := h.HydrateAggregate(ctx, c.AggregateType(), c.AggregateId())
	if err != nil {
		return new(T), fmt.Errorf("failed to hydrate aggregate(%s#%s): %w", c.AggregateType(), c.AggregateId(), err)
	}

	// check command validity for aggregate
	aggregate, events, err := h.ApplyCommand(ctx, aggregate, c)
	if err != nil {
		return new(T), fmt.Errorf("command (%T) rejected on aggregate(%s#%s): %w", c, c.AggregateType(), c.AggregateId(), err)
	}

	// persist and publish events
	err = h.PersistAndPublish(ctx, events...)
	if err != nil {
		return new(T), fmt.Errorf("failed to persist and publish events for aggregate(%s#%s): %w", c.AggregateType(), c.AggregateId(), err)
	}

	// return aggregate
	return aggregate, nil
}

func (h *commandHandler[T]) HydrateAggregate(ctx context.Context, aggregateType AggregateType, aggregateId uuid.UUID) (*T, error) {
	events, err := h.eventStore.Load(aggregateType, aggregateId)
	if err != nil {
		return new(T), fmt.Errorf("failed to load events for aggregate(%s#%s): %w", aggregateType, aggregateId, err)
	}

	// create new aggregate
	aggregate := h.factory()

	// apply events
	for _, event := range events {
		err := event.Apply(aggregate)
		if err != nil {
			return new(T), fmt.Errorf("failed to apply event(%s) to aggregate(%s#%s): %w", event.EventType(), aggregateType, aggregateId, err)
		}
	}

	// return aggregate
	return aggregate, nil
}

func (h *commandHandler[T]) ApplyCommand(ctx context.Context, aggregate *T, c Command[T]) (*T, []Event[T], error) {
	agg := *aggregate
	// check if command is valid for aggregate
	events, err := c.Apply(aggregate)
	if err != nil {
		return new(T), nil, fmt.Errorf("command (%T) is invalid for aggregate(%s#%s): %w", c, agg.AggregateType(), agg.AggregateId(), err)
	}

	for _, event := range events {
		err := event.Apply(aggregate)
		if err != nil {
			return new(T), nil, fmt.Errorf("failed to apply event(%s) to aggregate(%s#%s): %w", event.EventType(), agg.AggregateType(), agg.AggregateId(), err)
		}
	}

	// return events
	return aggregate, events, nil
}

// TODO: make it transactional
func (h *commandHandler[T]) PersistAndPublish(ctx context.Context, events ...Event[T]) error {
	err := h.eventStore.Store(ctx, events...)
	if err != nil {
		return fmt.Errorf("failed to store events: %w", err)
	}

	// publish events
	err = h.eventPublisher.Publish(events...)
	if err != nil {
		return fmt.Errorf("failed to publish events: %w", err)
	}

	return nil
}
