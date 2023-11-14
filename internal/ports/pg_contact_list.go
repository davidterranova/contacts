package ports

import (
	"context"
	"fmt"
	"time"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PgContactList struct {
	db *gorm.DB
}

type pgContact struct {
	Id               uuid.UUID  `gorm:"primaryKey;column:id"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at"`
	CreatedBy        string     `gorm:"column:created_by"`
	AggregateVersion int        `gorm:"column:aggregate_version"`

	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
	Phone     string `gorm:"column:phone"`
}

func (pgContact) TableName() string {
	return "contacts"
}

func NewPgContactList(db *gorm.DB, eventStream eventsourcing.Subscriber[domain.Contact]) *PgContactList {
	pgContactList := &PgContactList{
		db: db,
	}
	eventStream.Subscribe(context.Background(), pgContactList.HandleEvent)

	return pgContactList
}

func (l *PgContactList) HandleEvent(e eventsourcing.Event[domain.Contact]) {
	var err error

	log.Debug().
		Str("aggregate_id", e.AggregateId().String()).
		Str("aggregate_type", string(e.AggregateType())).
		Str("event_type", e.EventType()).
		Interface("event", e).
		Msg("read model handling event")

	switch e.EventType() {
	case domain.ContactCreated:
		err = l.create(pgContact{
			Id:               e.AggregateId(),
			CreatedAt:        e.IssuedAt(),
			UpdatedAt:        e.IssuedAt(),
			CreatedBy:        e.IssuedBy().String(),
			AggregateVersion: e.AggregateVersion(),
		})
	case domain.ContactEmailUpdated:
		ev, ok := e.(*domain.EvtContactEmailUpdated)
		if !ok {
			err = fmt.Errorf("%w: %T", ErrUnknownEvent, e)
			break
		}
		err = l.update(e.AggregateId(), func(pgC pgContact) pgContact {
			pgC.AggregateVersion = e.AggregateVersion()
			pgC.UpdatedAt = e.IssuedAt()
			pgC.Email = ev.Email
			return pgC
		})
	case domain.ContactNameUpdated:
		ev, ok := e.(*domain.EvtContactNameUpdated)
		if !ok {
			err = fmt.Errorf("%w: %T", ErrUnknownEvent, e)
			break
		}
		err = l.update(e.AggregateId(), func(pgC pgContact) pgContact {
			pgC.AggregateVersion = e.AggregateVersion()
			pgC.UpdatedAt = e.IssuedAt()
			pgC.FirstName = ev.FirstName
			pgC.LastName = ev.LastName
			return pgC
		})
	case domain.ContactPhoneUpdated:
		ev, ok := e.(*domain.EvtContactPhoneUpdated)
		if !ok {
			err = fmt.Errorf("%w: %T", ErrUnknownEvent, e)
			break
		}
		err = l.update(e.AggregateId(), func(pgC pgContact) pgContact {
			pgC.AggregateVersion = e.AggregateVersion()
			pgC.UpdatedAt = e.IssuedAt()
			pgC.Phone = ev.Phone
			return pgC
		})
	case domain.ContactDeleted:
		err = l.delete(e.AggregateId())
	default:
		err = fmt.Errorf("%w: %T", ErrUnknownEvent, e)
	}

	if err != nil {
		log.Error().Err(err).
			Interface("event", e).
			Str("aggregate_type", string(e.AggregateType())).
			Str("event_type", e.EventType()).
			Str("aggregate_id", e.AggregateId().String()).
			Msg("error applying event")
	}
}

func (l *PgContactList) List(ctx context.Context, query usecase.QueryListContact) ([]*domain.Contact, error) {
	var pgContacts []pgContact

	err := l.db.WithContext(ctx).
		Find(&pgContacts).
		Where("created_by = ?", query.User.Id().String()).
		Order("created_at DESC").
		Error
	if err != nil {
		return nil, err
	}

	contacts := make([]*domain.Contact, 0, len(pgContacts))
	for _, pgContact := range pgContacts {
		contacts = append(contacts, fromPgContact(pgContact))
	}

	return contacts, nil
}

func (l *PgContactList) create(c pgContact) error {
	return l.db.Create(&c).Error
}

func (l *PgContactList) load(id uuid.UUID) (pgContact, error) {
	var pgC pgContact
	err := l.db.Take(&pgC, id).Error

	return pgC, err
}

// func (l *PgContactList) update(c pgContact) error {
// 	return l.db.Save(&c).Error
// }

func (l *PgContactList) update(id uuid.UUID, fn func(pgC pgContact) pgContact) error {
	return l.db.Transaction(func(tx *gorm.DB) error {
		pgC, err := l.load(id)
		if err != nil {
			return err
		}

		pgC = fn(pgC)

		return tx.Save(&pgC).Error
	})
}

func (l *PgContactList) delete(id uuid.UUID) error {
	return l.db.Delete(&pgContact{}, id).Error
}

func fromPgContact(pg pgContact) *domain.Contact {
	return &domain.Contact{
		AggregateBase: eventsourcing.NewAggregateBase[domain.Contact](pg.Id),
		CreatedAt:     pg.CreatedAt,
		UpdatedAt:     pg.UpdatedAt,
		DeletedAt:     pg.DeletedAt,
		CreatedBy:     user.FromString(pg.CreatedBy),
		FirstName:     pg.FirstName,
		LastName:      pg.LastName,
		Email:         pg.Email,
		Phone:         pg.Phone,
	}
}
