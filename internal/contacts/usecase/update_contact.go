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

type CmdUpdateContact struct {
	ContactId string `validate:"required,uuid"`
	FirstName string `validate:"omitempty,min=2,max=255"`
	LastName  string `validate:"omitempty,min=2,max=255"`
	Email     string `validate:"omitempty,email"`
	Phone     string `validate:"omitempty,e164"` // https://en.wikipedia.org/wiki/E.164
}

type UpdateContact struct {
	validator      *validator.Validate
	commandHandler eventsourcing.CommandHandler[domain.Contact]
}

func NewUpdateContact(commandHandler eventsourcing.CommandHandler[domain.Contact]) UpdateContact {
	return UpdateContact{
		validator:      validator.New(),
		commandHandler: commandHandler,
	}
}

func (h UpdateContact) Update(ctx context.Context, cmd CmdUpdateContact, cmdIssuedBy user.User) (*domain.Contact, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	uuid, err := uuid.Parse(cmd.ContactId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	checkedCmd := newCmdUpdateContact(uuid, cmd, cmdIssuedBy)
	return handleErrs(h.commandHandler.HandleCommand(ctx, checkedCmd))
}

type cmdUpdateContact struct {
	eventsourcing.BaseCommand[*domain.Contact]
	CmdUpdateContact
}

func newCmdUpdateContact(contactId uuid.UUID, data CmdUpdateContact, cmdIssuedBy user.User) cmdUpdateContact {
	return cmdUpdateContact{
		BaseCommand: eventsourcing.NewBaseCommand[*domain.Contact](
			contactId,
			domain.AggregateContact,
			cmdIssuedBy,
		),
		CmdUpdateContact: data,
	}
}

func (c cmdUpdateContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[domain.Contact], error) {
	err := eventsourcing.EnsureAggregateNotNew(aggregate)
	if err != nil {
		return nil, err
	}

	if err := checkUpdatePolicy(c, aggregate); err != nil {
		return nil, err
	}
	if aggregate.DeletedAt != nil {
		return nil, ErrNotFound
	}

	aggregateVersion := aggregate.AggregateVersion()
	events := make([]eventsourcing.Event[domain.Contact], 0)
	if c.FirstName != "" || c.LastName != "" {
		updateNames := false
		firstName := aggregate.FirstName
		lastName := aggregate.LastName

		if c.FirstName != "" && c.FirstName != aggregate.FirstName {
			firstName = c.FirstName
			updateNames = true
		}
		if c.LastName != "" && c.LastName != aggregate.LastName {
			lastName = c.LastName
			updateNames = true
		}

		if updateNames {
			aggregateVersion++
			events = append(events, domain.NewEvtContactNameUpdated(c.AggregateId(), aggregateVersion, c.IssuedBy(), firstName, lastName))
		}
	}

	if c.Email != "" && aggregate.Email != c.Email {
		aggregateVersion++
		events = append(events, domain.NewEvtContactEmailUpdated(c.AggregateId(), aggregateVersion, c.IssuedBy(), c.Email))
	}
	if c.Phone != "" && aggregate.Phone != c.Phone {
		aggregateVersion++
		events = append(events, domain.NewEvtContactPhoneUpdated(c.AggregateId(), aggregateVersion, c.IssuedBy(), c.Phone))
	}

	return events, nil
}

func checkUpdatePolicy(cmd cmdUpdateContact, aggregate *domain.Contact) error {
	if cmd.IssuedBy().Id() == aggregate.CreatedBy.Id() {
		return nil
	}

	return ErrForbidden
}
