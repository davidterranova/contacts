package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdDeleteContact struct {
	eventsourcing.BaseCommand[*domain.Contact] `validate:"required"`
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

func (h DeleteContactHandler) Delete(ctx context.Context, cmd CmdDeleteContact) error {
	err := h.validator.Struct(cmd)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	_, err = h.commandHandler.Handle(cmd)
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

func NewCmdDeleteContact(contactId uuid.UUID) CmdDeleteContact {
	return CmdDeleteContact{
		BaseCommand: eventsourcing.NewBaseCommand[*domain.Contact](
			contactId,
			domain.AggregateContact,
		),
	}
}

func (c CmdDeleteContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[*domain.Contact], error) {
	if aggregate.AggregateId() == uuid.Nil {
		return nil, eventsourcing.ErrAggregateNotFound
	}
	if aggregate.DeletedAt != nil {
		return nil, ErrNotFound
	}

	return []eventsourcing.Event[*domain.Contact]{
		domain.NewEvtContactDeleted(c.AggregateId()),
	}, nil
}
