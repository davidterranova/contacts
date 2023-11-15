package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davidterranova/contacts/internal/admin/domain"
	"github.com/davidterranova/contacts/internal/admin/usecase"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PgListEvent struct {
	db *gorm.DB
}

type pgEvent struct {
	EventId          uuid.UUID                   `gorm:"type:uuid;primaryKey;column:event_id"`
	EventType        string                      `gorm:"type:varchar(255);column:event_type"`
	EventIssuedAt    time.Time                   `gorm:"column:event_issued_at"`
	EventIssuedBy    string                      `gorm:"type:varchar(255);column:event_issued_by"`
	EventData        json.RawMessage             `gorm:"type:jsonb;column:event_data"`
	AggregateId      uuid.UUID                   `gorm:"type:uuid;column:aggregate_id"`
	AggregateType    eventsourcing.AggregateType `gorm:"type:varchar(255);column:aggregate_type"`
	AggregateVersion int                         `gorm:"column:aggregate_version"`
	Published        bool                        `gorm:"column:published"`
}

func (pgEvent) TableName() string {
	return "events"
}

func NewPgListEvent(db *gorm.DB) *PgListEvent {
	return &PgListEvent{
		db: db,
	}
}

func AggregateIdScope(aggregateId *uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if aggregateId == nil || *aggregateId == uuid.Nil {
			return db
		}

		return db.Where("aggregate_id = ?", aggregateId)
	}
}

func AggregateTypeScope(aggregateType *eventsourcing.AggregateType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if aggregateType == nil || *aggregateType == "" {
			return db
		}

		return db.Where("aggregate_type = ?", aggregateType)
	}
}

func AggregateTypePublished(published *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if published == nil {
			return db
		}

		return db.Where("published = ?", published)
	}
}

func (l *PgListEvent) ListEvents(ctx context.Context, query usecase.QueryListEvent) ([]*domain.Event, error) {
	var pgEvents []pgEvent

	dbQuery := l.db.WithContext(ctx).
		Table(pgEvent{}.TableName()).
		Joins("JOIN events_outbox ON events_outbox.event_id = events.event_id")
	dbQuery = dbQuery.Scopes(
		AggregateIdScope(query.AggregateId),
		AggregateTypeScope(query.AggregateType),
		AggregateTypePublished(query.Published),
	)

	// err := dbQuery.Find(&pgEvents).Error
	err := dbQuery.Scan(&pgEvents).Error
	if err != nil {
		return nil, fmt.Errorf("error listing events: %w", err)
	}

	events := make([]*domain.Event, 0, len(pgEvents))
	for _, pgEvent := range pgEvents {
		event, err := toEvent(pgEvent)
		if err != nil {
			return nil, fmt.Errorf("error listing events: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func toEvent(e pgEvent) (*domain.Event, error) {
	var u user.User
	err := json.Unmarshal([]byte(e.EventIssuedBy), &u)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal user", err)
	}

	return &domain.Event{
		EventId:          e.EventId,
		EventType:        e.EventType,
		EventIssuesAt:    e.EventIssuedAt,
		EventIssuedBy:    u,
		EventData:        e.EventData,
		AggregateId:      e.AggregateId,
		AggregateType:    e.AggregateType,
		AggregateVersion: e.AggregateVersion,
		Published:        e.Published,
	}, nil
}
