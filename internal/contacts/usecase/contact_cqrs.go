//go:generate mockgen -destination=mock_contact_cmd_handler.go -package=usecase . ContactCmdHandler
package usecase

import (
	"context"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/cqrs/eventsourcing"
	uuid "github.com/google/uuid"
	_ "go.uber.org/mock/mockgen/model"
)

// ContactCmdHandler is a mock of eventsourcing.CommandHandler interface.
type ContactCmdHandler interface {
	// Handle is the global command handler that should be called by the application
	// HandleCommand(ctx context.Context, cmd eventsourcing.Command[domain.Contact]) (*domain.Contact, error)

	// HydrateAggregate an aggregate from already published events (internal)
	// HydrateAggregate(ctx context.Context, aggregateType eventsourcing.AggregateType, aggregateId uuid.UUID) (*domain.Contact, error)

	// Apply checks command validity for an aggregate and return newly emitted events (internal)
	// ApplyCommand(ctx context.Context, aggregate *domain.Contact, command eventsourcing.Command[domain.Contact]) (*domain.Contact, []eventsourcing.Event[domain.Contact], error)

	// HandleCommand is the global command handler that should be called by the application
	HandleCommand(ctx context.Context, cmd eventsourcing.Command[domain.Contact]) (*domain.Contact, error)

	// HydrateAggregate an aggregate from all its events (internal)
	HydrateAggregate(ctx context.Context, aggregateType eventsourcing.AggregateType, aggregateId uuid.UUID) (*domain.Contact, error)
}

type ContactLister interface {
	List(ctx context.Context, query QueryListContact) ([]*domain.Contact, error)
}
