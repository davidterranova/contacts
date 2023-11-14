package eventsourcing

import "github.com/google/uuid"

type AggregateType string

type Aggregate interface {
	AggregateId() uuid.UUID
	// SetAggregateId(uuid.UUID)
	AggregateType() AggregateType

	IncrementVersion()
	AggregateVersion() int
}

type AggregateBase struct {
	aggregateId      uuid.UUID
	aggregateVersion int
}

func NewAggregateBase(aggregateId uuid.UUID) *AggregateBase {
	return &AggregateBase{
		aggregateId:      aggregateId,
		aggregateVersion: 0,
	}
}

func (a AggregateBase) AggregateId() uuid.UUID {
	return a.aggregateId
}

func (a *AggregateBase) SetAggregateId(aggregateId uuid.UUID) {
	a.aggregateId = aggregateId
}

func (a *AggregateBase) IncrementVersion() {
	a.aggregateVersion++
}

func (a AggregateBase) AggregateVersion() int {
	return a.aggregateVersion
}
