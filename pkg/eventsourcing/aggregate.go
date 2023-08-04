package eventsourcing

import "github.com/google/uuid"

type AggregateType string

type Aggregate interface {
	AggregateId() uuid.UUID
	AggregateType() AggregateType
}
