package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdDeleteContact struct {
	ContactId string `validate:"required,uuid"`
}

type DeleteContactHandler struct {
	validator      *validator.Validate
	commandHandler eventsourcing.CommandHandler[*domain.Contact]
}

func NewDeleteContact(commandHandler eventsourcing.CommandHandler[*domain.Contact]) DeleteContactHandler {
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

	_, err = h.commandHandler.Handle(newCmdDeleteContact(uuid, cmdIssuedBy))
	switch {
	case errors.Is(err, ErrNotFound):
		return err
	case errors.Is(err, eventsourcing.ErrAggregateNotFound):
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	case err != nil:
		return fmt.Errorf("%w: %s", ErrInternal, err)
	default:
		return nil
	}
}

type cmdDeleteContact struct {
	eventsourcing.BaseCommand[*domain.Contact]
}

func newCmdDeleteContact(contactId uuid.UUID, cmdIssuedBy user.User) cmdDeleteContact {
	return cmdDeleteContact{
		BaseCommand: eventsourcing.NewBaseCommand[*domain.Contact](
			contactId,
			domain.AggregateContact,
			cmdIssuedBy,
		),
	}
}

func (c cmdDeleteContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[*domain.Contact], error) {
	if aggregate.AggregateId() == uuid.Nil {
		return nil, eventsourcing.ErrAggregateNotFound
	}
	if err := checkDeletePolicy(c, aggregate); err != nil {
		return nil, err
	}
	if aggregate.DeletedAt != nil {
		return nil, ErrNotFound
	}

	return []eventsourcing.Event[*domain.Contact]{
		domain.NewEvtContactDeleted(
			c.AggregateId(),
			c.IssuedBy(),
		),
	}, nil
}

func checkDeletePolicy(cmd cmdDeleteContact, aggregate *domain.Contact) error {
	if cmd.IssuedBy().Id() == aggregate.CreatedBy.Id() {
		return nil
	}

	return ErrForbidden
}
