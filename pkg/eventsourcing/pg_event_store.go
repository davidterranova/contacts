package eventsourcing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

	AggregateId      uuid.UUID     `gorm:"type:uuid;column:aggregate_id"`
	AggregateType    AggregateType `gorm:"type:varchar(255);column:aggregate_type"`
	AggregateVersion int           `gorm:"column:aggregate_version"`
}

func (pgEvent) TableName() string {
	return "events"
}

type pgOutboxEntry struct {
	EventId          uuid.UUID `gorm:"type:uuid;primaryKey;column:event_id"`
	Published        bool      `gorm:"column:published"`
	AggregateVersion int       `gorm:"column:aggregate_version"`
}

func (pgOutboxEntry) TableName() string {
	return "events_outbox"
}

func NewPGEventStore[T Aggregate](db *gorm.DB, registry *Registry[T]) *pgEventStore[T] {
	return &pgEventStore[T]{
		db:       db,
		registry: registry,
	}
}

func (s *pgEventStore[T]) Store(ctx context.Context, events ...Event[T]) error {
	pgEvents := make([]*pgEvent, 0, len(events))
	outboxEntries := make([]*pgOutboxEntry, 0, len(events))

	for _, event := range events {
		pgEvent, err := s.toPgEvent(event)
		if err != nil {
			return err
		}
		pgEvents = append(pgEvents, pgEvent)

		outboxEntries = append(outboxEntries, &pgOutboxEntry{
			EventId:          event.Id(),
			Published:        false,
			AggregateVersion: event.AggregateVersion(),
		})
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(pgEvents).Error
		if err != nil {
			return fmt.Errorf("failed to create events in event_store table: %w", err)
		}

		for _, event := range events {
			log.Debug().Str("type", event.EventType()).Interface("event", event).Msg("stored event")
		}

		return tx.Create(outboxEntries).Error
	})
}

func (s *pgEventStore[T]) Load(ctx context.Context, aggregateType AggregateType, aggregateId uuid.UUID) ([]Event[T], error) {
	var pgEvents []pgEvent
	err := s.db.WithContext(ctx).Where("aggregate_type = ? AND aggregate_id = ?", aggregateType, aggregateId).Find(&pgEvents).Error
	if err != nil {
		return nil, err
	}

	return s.fromPgEvenSlice(pgEvents)
}

func (s *pgEventStore[T]) LoadUnpublished(ctx context.Context, batchSize int) ([]Event[T], error) {
	var pgOutboxEntries []uuid.UUID
	// err := s.db.WithContext(ctx).Where("published = ?", false).Limit(batchSize).Find(&pgOutboxEntries).Error
	err := s.db.
		WithContext(ctx).
		Model(&pgOutboxEntry{}).
		Where("published = ?", false).
		Group("event_id").
		Order("aggregate_version ASC").
		Limit(batchSize).
		Pluck("event_id", &pgOutboxEntries).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to load unpublished events from outbox: %w", err)
	}

	var unpublishedEvents []pgEvent
	err = s.db.WithContext(ctx).Where("event_id IN ?", pgOutboxEntries).Find(&unpublishedEvents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load unpublished events: %w", err)
	}

	for _, event := range unpublishedEvents {
		log.Debug().Str("type", event.EventType).Interface("event", event).Msg("loaded unpublished event")
	}

	return s.fromPgEvenSlice(unpublishedEvents)
}

func (s *pgEventStore[T]) MarkPublished(ctx context.Context, events ...Event[T]) error {
	var eventIds []uuid.UUID
	for _, event := range events {
		eventIds = append(eventIds, event.Id())
		log.Debug().Str("type", event.EventType()).Interface("event", event).Msg("marked event as published")
	}

	return s.db.WithContext(ctx).Model(&pgOutboxEntry{}).Where("event_id IN ?", eventIds).Update("published", true).Error
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
		EventId:          e.Id(),
		EventType:        e.EventType(),
		EventIssuedAt:    e.IssuedAt(),
		EventIssuedBy:    string(byteUser),
		EventData:        data,
		AggregateId:      e.AggregateId(),
		AggregateType:    e.AggregateType(),
		AggregateVersion: e.AggregateVersion(),
	}, nil
}

func (s *pgEventStore[T]) fromPgEvenSlice(pgEvents []pgEvent) ([]Event[T], error) {
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

func (s *pgEventStore[T]) fromPgEvent(pgEvent pgEvent) (Event[T], error) {
	var u user.User
	err := json.Unmarshal([]byte(pgEvent.EventIssuedBy), &u)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal user", err)
	}

	return s.registry.Hydrate(
		EventBase[T]{
			eventId:          pgEvent.EventId,
			eventIssuesAt:    pgEvent.EventIssuedAt,
			eventIssuedBy:    u,
			eventType:        pgEvent.EventType,
			aggregateType:    pgEvent.AggregateType,
			aggregateId:      pgEvent.AggregateId,
			aggregateVersion: pgEvent.AggregateVersion,
		},
		pgEvent.EventData,
	)
}
