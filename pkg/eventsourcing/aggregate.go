package eventsourcing

import "github.com/google/uuid"

type AggregateType string

type Aggregate interface {
	AggregateId() uuid.UUID
	AggregateType() AggregateType

	IncrementVersion()
	AggregateVersion() int
}

type AggregateBase struct {
	aggregateVersion int
}

func (a *AggregateBase) IncrementVersion() {
	a.aggregateVersion++
}

func (a AggregateBase) AggregateVersion() int {
	return a.aggregateVersion
}
