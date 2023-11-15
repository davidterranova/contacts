package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/admin/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/google/uuid"
)

type QueryListEvent struct {
	AggregateType *eventsourcing.AggregateType
	AggregateId   *uuid.UUID
	Published     *bool
}

type AggregateRepository interface {
	ListEvents(ctx context.Context, query QueryListEvent) ([]*domain.Event, error)
}
