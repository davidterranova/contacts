package usecase

import (
	"context"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/go-playground/validator"
	uuid "github.com/google/uuid"
)

type CmdCreateContact struct {
	FirstName string `validate:"min=2,max=255"`
	LastName  string `validate:"min=2,max=255"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"e164"` // https://en.wikipedia.org/wiki/E.164
}

type CreateContact struct {
	validator      *validator.Validate
	commandHandler eventsourcing.CommandHandler[domain.Contact]
}

func NewCreateContact(commandHandler eventsourcing.CommandHandler[domain.Contact]) CreateContact {
	return CreateContact{
		commandHandler: commandHandler,
		validator:      validator.New(),
	}
}

func (h CreateContact) Create(ctx context.Context, cmd CmdCreateContact, cmdIssuedBy user.User) (*domain.Contact, error) {
	err := h.validator.Struct(cmd)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, err)
	}

	checkedCmd := newCmdCreateContact(cmd, cmdIssuedBy)
	return handleErrs(h.commandHandler.Handle(checkedCmd))
}

type cmdCreateContact struct {
	eventsourcing.BaseCommand[*domain.Contact]
	CmdCreateContact
}

func newCmdCreateContact(data CmdCreateContact, cmdIssuedBy user.User) cmdCreateContact {
	return cmdCreateContact{
		BaseCommand: eventsourcing.NewBaseCommand[*domain.Contact](
			uuid.New(),
			domain.AggregateContact,
			cmdIssuedBy,
		),
		CmdCreateContact: data,
	}
}

func (c cmdCreateContact) Apply(aggregate *domain.Contact) ([]eventsourcing.Event[domain.Contact], error) {
	if aggregate.AggregateId() != uuid.Nil {
		return nil, eventsourcing.ErrAggregateAlreadyExists
	}

	return []eventsourcing.Event[domain.Contact]{
		domain.NewEvtContactCreated(c.AggregateId(), c.IssuedBy()),
		domain.NewEvtContactEmailUpdated(c.AggregateId(), c.IssuedBy(), c.Email),
		domain.NewEvtContactNameUpdated(c.AggregateId(), c.IssuedBy(), c.FirstName, c.LastName),
		domain.NewEvtContactPhoneUpdated(c.AggregateId(), c.IssuedBy(), c.Phone),
	}, nil
}
