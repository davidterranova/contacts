package internal

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
)

type ListContact interface {
	List(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error)
}

type CreateContact interface {
	Create(ctx context.Context, cmd usecase.CmdCreateContact) (*domain.Contact, error)
}

type UpdateContact interface {
	Update(ctx context.Context, cmd usecase.CmdUpdateContact) (*domain.Contact, error)
}

type DeleteContact interface {
	Delete(ctx context.Context, cmd usecase.CmdDeleteContact) error
}

type App struct {
	listContact   ListContact
	createContact CreateContact
	updateContact UpdateContact
	deleteContact DeleteContact
}

func New(repo usecase.ContactRepository) *App {
	return &App{
		listContact:   usecase.NewListContact(repo),
		createContact: usecase.NewCreateContact(repo),
		updateContact: usecase.NewUpdateContact(repo),
		deleteContact: usecase.NewDeleteContact(repo),
	}
}

func (a *App) ListContacts(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error) {
	return a.listContact.List(ctx, query)
}

func (a *App) CreateContact(ctx context.Context, cmd usecase.CmdCreateContact) (*domain.Contact, error) {
	return a.createContact.Create(ctx, cmd)
}

func (a *App) UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact) (*domain.Contact, error) {
	return a.updateContact.Update(ctx, cmd)
}

func (a *App) DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact) error {
	return a.deleteContact.Delete(ctx, cmd)
}
