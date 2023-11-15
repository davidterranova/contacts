package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/admin/domain"
)

type ListEventHandler struct {
	lister AggregateRepository
}

func NewListEvent(lister AggregateRepository) ListEventHandler {
	return ListEventHandler{
		lister: lister,
	}
}

func (h ListEventHandler) ListEvents(ctx context.Context, query QueryListEvent) ([]*domain.Event, error) {
	return h.lister.ListEvents(ctx, query)
}
