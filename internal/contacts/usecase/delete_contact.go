package usecase

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/cqrs/eventsourcing"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdDeleteContact struct {
	ContactId string `validate:"required,uuid"`
}

type DeleteContactHandler struct {
	validator      *validator.Validate
	commandHandler eventsourcing.CommandHandler[domain.Contact]
}

func NewDeleteContact(commandHandler eventsourcing.CommandHandler[domain.Contact]) DeleteContactHandler {
	return DeleteContactHandler{
		validator:      validator.New(),
		commandHandler: commandHandler,
	}
}

func (h DeleteContactHandler) Delete(ctx context.Context, cmd CmdDeleteContact, cmdIssuedBy user.User) error {
	err := h.validator.Struct(cmd)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	uuid, err := uuid.Parse(cmd.ContactId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	_, err = handleErrs(
		h.commandHandler.HandleCommand(
			ctx,
			newCmdDeleteContact(uuid, cmdIssuedBy),
		),
	)

	return err
}

type cmdDeleteContact struct {
	eventsourcing.BaseCommand[domain.Contact]
}

func newCmdDeleteContact(contactId uuid.UUID, cmdIssuedBy user.User) cmdDeleteContact {
	return cmdDeleteContact{
		BaseCommand: eventsourcing.NewBaseCommand[domain.Contact](
			contactId,
			domain.AggregateContact,
			cmdIssuedBy,
		),
	}
}

func (c cmdDeleteContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[domain.Contact], error) {
	err := eventsourcing.EnsureAggregateNotNew(aggregate)
	if err != nil {
		return nil, err
	}

	if err := checkDeletePolicy(c, aggregate); err != nil {
		return nil, err
	}
	if aggregate.DeletedAt != nil {
		return nil, ErrNotFound
	}

	return []eventsourcing.Event[domain.Contact]{
		domain.NewEvtContactDeleted(c.AggregateId(), aggregate.AggregateVersion()+1, c.IssuedBy()),
	}, nil
}

func checkDeletePolicy(cmd cmdDeleteContact, aggregate *domain.Contact) error {
	if cmd.IssuedBy().Id() == aggregate.CreatedBy.Id() {
		return nil
	}

	return ErrForbidden
}
