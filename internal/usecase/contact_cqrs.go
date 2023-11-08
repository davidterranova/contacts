//go:generate mockgen -destination=mock_contact_cmd_handler.go -package=usecase . ContactCmdHandler
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	_ "go.uber.org/mock/mockgen/model"
)

type ContactCmdHandler interface {
	Handle(eventsourcing.Command[*domain.Contact]) (*domain.Contact, error)
}

type ContactLister interface {
	List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error)
}
