package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/ports"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateContac(t *testing.T) {
	t.Parallel()

	testUpdateContactValidation(t)
	testUpdateContact(t)
	testUpdateContactFn(t)
}

func testUpdateContactValidation(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactUpdater := NewUpdateContact(container.contactRepo)

	testCases := []struct {
		name          string
		command       CmdUpdateContact
		expectedError error
	}{
		{
			name: "valid command",
			command: CmdUpdateContact{
				Updater:   *domain.NewUser(uuid.New()),
				ContactId: uuid.NewString(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "test@contact.local",
				Phone:     "+33612345678",
			},
			expectedError: nil,
		},
		{
			name: "invalid command: invalid contact id",
			command: CmdUpdateContact{
				ContactId: "invalid-uuid",
				FirstName: "John",
			},
			expectedError: ErrInvalidCommand,
		},
		{
			name: "invalid command: missing contact id",
			command: CmdUpdateContact{
				FirstName: "John",
			},
			expectedError: ErrInvalidCommand,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.expectedError == nil {
				container.contactRepo.EXPECT().
					Update(ctx, gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
			}

			_, err := contactUpdater.Update(ctx, tc.command)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testUpdateContact(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactUpdater := NewUpdateContact(container.contactRepo)

	t.Run("successfully update contact", func(t *testing.T) {
		uuid := uuid.New()
		cmd := CmdUpdateContact{
			Updater:   *domain.NewUser(uuid),
			ContactId: uuid.String(),
			FirstName: "John",
		}

		container.contactRepo.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				&domain.Contact{
					Id:        uuid,
					FirstName: cmd.FirstName,
				},
				nil,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd)
		assert.NoError(t, err)
		assert.Equal(t, cmd.ContactId, updatedContact.Id.String())
		assert.Equal(t, cmd.FirstName, updatedContact.FirstName)
	})

	t.Run("contact not found", func(t *testing.T) {
		cmd := CmdUpdateContact{
			Updater:   *domain.NewUser(uuid.New()),
			ContactId: uuid.NewString(),
			FirstName: "John",
		}

		container.contactRepo.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				nil,
				ports.ErrNotFound,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, updatedContact)
	})

	t.Run("unexpected repository error", func(t *testing.T) {
		cmd := CmdUpdateContact{
			Updater:   *domain.NewUser(uuid.New()),
			ContactId: uuid.NewString(),
			FirstName: "John",
		}

		container.contactRepo.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			Times(1).
			Return(
				nil,
				errors.New("internal error"),
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd)
		assert.ErrorIs(t, err, ErrInternal)
		assert.Nil(t, updatedContact)
	})
}

func testUpdateContactFn(t *testing.T) {
	contact := domain.Contact{
		Id:        uuid.New(),
		CreatedBy: uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jdoe@contact.local",
		Phone:     "+33612345678",
	}

	tests := []struct {
		name          string
		contact       domain.Contact
		cmd           CmdUpdateContact
		expectedError error
	}{
		{
			name:    "successfully update contact",
			contact: contact,
			cmd: CmdUpdateContact{
				Updater:   *domain.NewUser(contact.CreatedBy),
				FirstName: "Jane",
			},
			expectedError: nil,
		},
		{
			name:    "unauthorized updater",
			contact: contact,
			cmd: CmdUpdateContact{
				Updater:   *domain.NewUser(uuid.New()),
				FirstName: "Jane",
			},
			expectedError: ErrForbidden,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := updateContactFn(test.contact, test.cmd)
			assert.ErrorIs(t, err, test.expectedError)
		})
	}
}
