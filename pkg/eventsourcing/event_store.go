package eventsourcing

import (
	"context"

	"github.com/google/uuid"
)

type EventStore[T Aggregate] interface {
	// Store events
	Store(ctx context.Context, events ...Event[T]) error
	// Load events from the given aggregate
	Load(ctx context.Context, aggregateType AggregateType, aggregateId uuid.UUID) ([]Event[T], error) // TODO: remove aggregateType

	// LoadUnpublished loads a batch of un published events
	LoadUnpublished(ctx context.Context, batchSize int) ([]Event[T], error)
	// MarkPublished marks events as published
	MarkPublished(ctx context.Context, events ...Event[T]) error
}
