package contacts

import (
	"context"
	"io"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/cqrs/eventsourcing"
)

type ContactWriteModel interface {
	eventsourcing.CommandHandler[domain.Contact]
}

type ContactReadModel interface {
	usecase.ContactReadModel
}

type ListContact interface {
	List(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error)
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

type ExportContact interface {
	Export(ctx context.Context, cmd usecase.CmdExportContact, cmdIssuedBy user.User) (io.Writer, error)
}

type App struct {
	listContact   ListContact
	createContact CreateContact
	updateContact UpdateContact
	deleteContact DeleteContact
	exportContact ExportContact
}

func New(writeModel ContactWriteModel, readModel ContactReadModel) *App {
	return &App{
		listContact:   usecase.NewListContact(readModel),
		exportContact: usecase.NewExportContact(readModel),
		createContact: usecase.NewCreateContact(writeModel),
		updateContact: usecase.NewUpdateContact(writeModel),
		deleteContact: usecase.NewDeleteContact(writeModel),
	}
}

func (a *App) ListContacts(ctx context.Context, cmdIssuedBy user.User) ([]*domain.Contact, error) {
	return a.listContact.List(ctx, cmdIssuedBy)
}

func (a *App) ExportContact(ctx context.Context, cmd usecase.CmdExportContact, cmdIssuedBy user.User) (io.Writer, error) {
	return a.exportContact.Export(ctx, cmd, cmdIssuedBy)
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
