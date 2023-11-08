package ports

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("contact not found")

// InMemoryContactRepository is not thread safe
type InMemoryContactRepository struct {
	contacts map[uuid.UUID]*domain.Contact
}

type filter struct {
	createdBy *uuid.UUID
}

func (f *filter) CreatedBy() *uuid.UUID {
	return f.createdBy
}

func WithCreatedBy(id uuid.UUID) withFilter {
	return func(f *filter) {
		f.createdBy = &id
	}
}

type withFilter func(f *filter)

func NewFilter(filters ...withFilter) *filter {
	filter := &filter{}
	for _, f := range filters {
		f(filter)
	}

	return filter
}

func NewInMemoryContactRepository() *InMemoryContactRepository {
	return &InMemoryContactRepository{
		contacts: map[uuid.UUID]*domain.Contact{},
	}
}

func (r *InMemoryContactRepository) Get(_ context.Context, id uuid.UUID) (*domain.Contact, error) {
	contact, ok := r.contacts[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	return contact, nil
}

func (r *InMemoryContactRepository) List(ctx context.Context, filter domain.Filter) ([]*domain.Contact, error) {
	contacts := make([]*domain.Contact, 0, len(r.contacts))
	for _, contact := range r.contacts {
		if filterBy(filter, contact) {
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}

func filterBy(filter domain.Filter, contact *domain.Contact) bool {
	if filter.CreatedBy() != nil && *filter.CreatedBy() != contact.CreatedBy {
		return false
	}

	return true
}

func (r *InMemoryContactRepository) Create(_ context.Context, contact *domain.Contact) (*domain.Contact, error) {
	return r.save(contact)
}

func (r *InMemoryContactRepository) Update(ctx context.Context, id uuid.UUID, updateFn func(c domain.Contact) (domain.Contact, error)) (*domain.Contact, error) {
	originalConact, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedContact, err := updateFn(*originalConact)
	if err != nil {
		return nil, err
	}

	return r.save(&updatedContact)
}

func (r *InMemoryContactRepository) save(contact *domain.Contact) (*domain.Contact, error) {
	r.contacts[contact.Id] = contact
	return contact, nil
}

func (r *InMemoryContactRepository) Delete(_ context.Context, id uuid.UUID, deleterFn func(c domain.Contact) error) error {
	contact, ok := r.contacts[id]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, id)
	}

	if err := deleterFn(*contact); err != nil {
		return err
	}

	delete(r.contacts, id)

	return nil
}
