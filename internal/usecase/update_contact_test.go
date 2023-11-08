package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/user"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateContac(t *testing.T) {
	t.Parallel()

	testUpdateContactValidation(t)
	testUpdateContact(t)
}

func testUpdateContactValidation(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactUpdater := NewUpdateContact(container.contactCmdHandler)
	cmdIssuer := user.New(uuid.New(), user.UserTypeAuthenticated)

	testCases := []struct {
		name          string
		command       CmdUpdateContact
		expectedError error
	}{
		{
			name: "valid command",
			command: CmdUpdateContact{
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
				container.contactCmdHandler.EXPECT().
					Handle(gomock.Any()).
					Times(1).
					Return(nil, nil)
			}

			_, err := contactUpdater.Update(ctx, tc.command, cmdIssuer)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testUpdateContact(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactUpdater := NewUpdateContact(container.contactCmdHandler)
	cmdIssuer := user.New(uuid.New(), user.UserTypeAuthenticated)

	t.Run("successfully update contact", func(t *testing.T) {
		uuid := uuid.New()
		cmd := CmdUpdateContact{
			ContactId: uuid.String(),
			FirstName: "John",
		}

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				&domain.Contact{
					Id:        uuid,
					FirstName: cmd.FirstName,
				},
				nil,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd, cmdIssuer)
		assert.NoError(t, err)
		assert.Equal(t, cmd.ContactId, updatedContact.Id.String())
		assert.Equal(t, cmd.FirstName, updatedContact.FirstName)
	})

	t.Run("contact not found", func(t *testing.T) {
		cmd := CmdUpdateContact{
			ContactId: uuid.NewString(),
			FirstName: "John",
		}

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				nil,
				eventsourcing.ErrAggregateNotFound,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd, cmdIssuer)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, updatedContact)
	})

	t.Run("unexpected repository error", func(t *testing.T) {
		cmd := CmdUpdateContact{
			ContactId: uuid.NewString(),
			FirstName: "John",
		}

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				nil,
				errors.New("internal error"),
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd, cmdIssuer)
		assert.ErrorIs(t, err, ErrInternal)
		assert.Nil(t, updatedContact)
	})

	t.Run("not owning contact", func(t *testing.T) {
		cmd := CmdUpdateContact{
			ContactId: uuid.NewString(),
			FirstName: "John",
		}
		cmdWrongIssuer := user.New(uuid.New(), user.UserTypeAuthenticated)

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				nil,
				ErrForbidden,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd, cmdWrongIssuer)
		assert.ErrorIs(t, err, ErrForbidden)
		assert.Nil(t, updatedContact)
	})
}
