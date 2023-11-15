package admin

import (
	"context"

	"github.com/davidterranova/contacts/internal/admin/domain"
	"github.com/davidterranova/contacts/internal/admin/usecase"
)

type ListEvent interface {
	ListEvents(ctx context.Context, query usecase.QueryListEvent) ([]*domain.Event, error)
}

type App struct {
	listEvent ListEvent
}

func New(listEvent ListEvent) *App {
	return &App{
		listEvent: listEvent,
	}
}

func (a *App) ListEvents(ctx context.Context, query usecase.QueryListEvent) ([]*domain.Event, error) {
	return a.listEvent.ListEvents(ctx, query)
}
