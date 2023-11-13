//go:generate mockgen -destination=mock_contact_cmd_handler.go -package=usecase . ContactCmdHandler
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	uuid "github.com/google/uuid"
	_ "go.uber.org/mock/mockgen/model"
)

// ContactCmdHandler is a mock of eventsourcing.CommandHandler interface.
type ContactCmdHandler interface {
	// Handle is the global command handler that should be called by the application
	Handle(eventsourcing.Command[domain.Contact]) (*domain.Contact, error)

	// HydrateAggregate an aggregate from already published events (internal)
	HydrateAggregate(aggregateType eventsourcing.AggregateType, aggregateId uuid.UUID) (*domain.Contact, error)

	// Apply checks command validity for an aggregate and return newly emitted events (internal)
	ApplyCommand(aggregate *domain.Contact, command eventsourcing.Command[domain.Contact]) (*domain.Contact, []eventsourcing.Event[domain.Contact], error)
}

type ContactLister interface {
	List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error)
}
