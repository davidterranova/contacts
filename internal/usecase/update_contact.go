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

type CmdUpdateContact struct {
	eventsourcing.BaseCommand[*domain.Contact] `validate:"required"`

	FirstName string `validate:"omitempty,min=2,max=255"`
	LastName  string `validate:"omitempty,min=2,max=255"`
	Email     string `validate:"omitempty,email"`
	Phone     string `validate:"omitempty,e164"` // https://en.wikipedia.org/wiki/E.164
}

type UpdateContact struct {
	validator      *validator.Validate
	commandHandler eventsourcing.CommandHandler[*domain.Contact]
}

func NewUpdateContact(commandHandler eventsourcing.CommandHandler[*domain.Contact]) UpdateContact {
	return UpdateContact{
		validator:      validator.New(),
		commandHandler: commandHandler,
	}
}

func (h UpdateContact) Update(ctx context.Context, cmd CmdUpdateContact) (*domain.Contact, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	contact, err := h.commandHandler.Handle(cmd)
	switch {
	case errors.Is(err, ErrNotFound):
		return nil, err
	case errors.Is(err, eventsourcing.ErrAggregateNotFound):
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	case err != nil:
		return nil, fmt.Errorf("%w: %s", ErrInternal, err)
	default:
		return contact, nil
	}
}

func NewCmdUpdateContact(contactId uuid.UUID, firstName, lastName, email, phone string) CmdUpdateContact {
	return CmdUpdateContact{
		BaseCommand: eventsourcing.NewBaseCommand[*domain.Contact](
			contactId,
			domain.AggregateContact,
		),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
}

func (c CmdUpdateContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[*domain.Contact], error) {
	if aggregate.AggregateId() == uuid.Nil {
		return nil, eventsourcing.ErrAggregateNotFound
	}
	if aggregate.DeletedAt != nil {
		return nil, ErrNotFound
	}

	events := make([]eventsourcing.Event[*domain.Contact], 0)
	if c.FirstName != "" || c.LastName != "" {
		if c.FirstName != "" && c.FirstName != aggregate.FirstName {
			aggregate.FirstName = c.FirstName
		}
		if c.LastName != "" && c.LastName != aggregate.LastName {
			aggregate.LastName = c.LastName
		}

		events = append(events, domain.NewEvtContactNameUpdated(c.AggregateId(), aggregate.FirstName, aggregate.LastName))
	}

	if c.Email != "" && aggregate.Email != c.Email {
		events = append(events, domain.NewEvtContactEmailUpdated(c.AggregateId(), c.Email))
	}
	if c.Phone != "" && aggregate.Phone != c.Phone {
		events = append(events, domain.NewEvtContactPhoneUpdated(c.AggregateId(), c.Phone))
	}

	return events, nil
}
