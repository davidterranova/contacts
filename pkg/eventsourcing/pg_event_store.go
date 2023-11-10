package eventsourcing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgEventStore[T Aggregate] struct {
	db       *gorm.DB
	registry *Registry[T]
}

type pgEvent struct {
	EventId       uuid.UUID       `gorm:"type:uuid;primaryKey;column:event_id"`
	EventType     string          `gorm:"type:varchar(255);column:event_type"`
	EventIssuedAt time.Time       `gorm:"column:event_issued_at"`
	EventIssuedBy string          `gorm:"type:varchar(255);column:event_issued_by"`
	EventData     json.RawMessage `gorm:"type:jsonb;column:event_data"`

	AggregateId   uuid.UUID     `gorm:"type:uuid;column:aggregate_id"`
	AggregateType AggregateType `gorm:"type:varchar(255);column:aggregate_type"`
}

func (pgEvent) TableName() string {
	return "events"
}

func NewPGEventStore[T Aggregate](db *gorm.DB, registry *Registry[T]) *pgEventStore[T] {
	return &pgEventStore[T]{
		db:       db,
		registry: registry,
	}
}

func (s *pgEventStore[T]) Store(events ...Event[T]) error {
	pgEvents := make([]*pgEvent, 0, len(events))
	for _, event := range events {
		pgEvent, err := s.toPgEvent(event)
		if err != nil {
			return err
		}
		pgEvents = append(pgEvents, pgEvent)
	}

	return s.db.Create(pgEvents).Error
}

func (s *pgEventStore[T]) Load(aggregateType AggregateType, aggregateId uuid.UUID) ([]Event[T], error) {
	var pgEvents []pgEvent
	err := s.db.Where("aggregate_type = ? AND aggregate_id = ?", aggregateType, aggregateId).Find(&pgEvents).Error
	if err != nil {
		return nil, err
	}

	events := make([]Event[T], 0, len(pgEvents))
	for _, pgEvent := range pgEvents {
		hydratedEvent, err := s.fromPgEvent(pgEvent)
		if err != nil {
			return nil, err
		}

		events = append(events, hydratedEvent)
	}

	return events, nil
}

func (s *pgEventStore[T]) toPgEvent(e Event[T]) (*pgEvent, error) {
	byteUser, err := e.IssuedBy().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal user", err)
	}

	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal event", err)
	}

	return &pgEvent{
		EventId:       e.Id(),
		EventType:     e.EventType(),
		EventIssuedAt: e.IssuedAt(),
		EventIssuedBy: string(byteUser),
		EventData:     data,
		AggregateId:   e.AggregateId(),
		AggregateType: e.AggregateType(),
	}, nil
}

func (s *pgEventStore[T]) fromPgEvent(pgEvent pgEvent) (Event[T], error) {
	var u user.User
	err := json.Unmarshal([]byte(pgEvent.EventIssuedBy), &u)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal user", err)
	}

	return s.registry.Hydrate(
		EventBase[T]{
			eventId:       pgEvent.EventId,
			eventIssuesAt: pgEvent.EventIssuedAt,
			eventIssuedBy: u,
			eventType:     pgEvent.EventType,
			aggregateType: pgEvent.AggregateType,
			aggregateId:   pgEvent.AggregateId,
		},
		pgEvent.EventData,
	)
}
