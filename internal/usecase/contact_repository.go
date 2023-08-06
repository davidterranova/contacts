//go:generate mockgen -destination=mock_contact_cmd_handler.go -package=usecase . ContactCmdHandler
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	_ "github.com/golang/mock/mockgen/model"
)

type ContactCmdHandler interface {
	Handle(eventsourcing.Command[*domain.Contact]) (*domain.Contact, error)
}

type ContactLister interface {
	List(ctx context.Context) ([]*domain.Contact, error)
}
