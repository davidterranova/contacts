package contacts

import (
	"context"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
)

type ContactWriteModel interface {
	eventsourcing.CommandHandler[domain.Contact]
}

type ContactReadModel interface {
	usecase.ContactLister
}

type ListContact interface {
	List(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error)
}

type CreateContact interface {
	Create(ctx context.Context, cmd usecase.CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error)
}

type UpdateContact interface {
	Update(ctx context.Context, cmd usecase.CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error)
}

type DeleteContact interface {
	Delete(ctx context.Context, cmd usecase.CmdDeleteContact, cmdIssuedBy user.User) error
}

type App struct {
	listContact   ListContact
	createContact CreateContact
	updateContact UpdateContact
	deleteContact DeleteContact
}

func New(writeModel ContactWriteModel, readModel ContactReadModel) *App {
	return &App{
		listContact:   usecase.NewListContact(readModel),
		createContact: usecase.NewCreateContact(writeModel),
		updateContact: usecase.NewUpdateContact(writeModel),
		deleteContact: usecase.NewDeleteContact(writeModel),
	}
}

func (a *App) ListContacts(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error) {
	return a.listContact.List(ctx, query)
}

func (a *App) CreateContact(ctx context.Context, cmd usecase.CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error) {
	return a.createContact.Create(ctx, cmd, cmdIssuedBy)
}

func (a *App) UpdateContact(ctx context.Context, cmd usecase.CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error) {
	return a.updateContact.Update(ctx, cmd, cmdIssuedBy)
}

func (a *App) DeleteContact(ctx context.Context, cmd usecase.CmdDeleteContact, cmdIssuedBy user.User) error {
	return a.deleteContact.Delete(ctx, cmd, cmdIssuedBy)
}
