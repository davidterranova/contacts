package eventsourcing

import "github.com/google/uuid"

type AggregateType string

type Aggregate interface {
	AggregateId() uuid.UUID
	AggregateType() AggregateType

	IncrementVersion()
	AggregateVersion() int
}

type AggregateBase[T Aggregate] struct {
	aggregateId      uuid.UUID
	aggregateVersion int
	events           []Event[T]
}

func NewAggregateBase[T Aggregate](aggregateId uuid.UUID) *AggregateBase[T] {
	return &AggregateBase[T]{
		aggregateId:      aggregateId,
		aggregateVersion: 0,
		events:           make([]Event[T], 0),
	}
}

func (a AggregateBase[T]) AggregateId() uuid.UUID {
	return a.aggregateId
}

func (a *AggregateBase[T]) Process(e Event[T]) {
	a.aggregateId = e.AggregateId()
	a.aggregateVersion = e.AggregateVersion()
	a.events = append(a.events, e)
}

func (a *AggregateBase[T]) IncrementVersion() {
	a.aggregateVersion++
}

func (a AggregateBase[T]) AggregateVersion() int {
	return a.aggregateVersion
}

func (a AggregateBase[T]) Events() []Event[T] {
	return a.events
}
