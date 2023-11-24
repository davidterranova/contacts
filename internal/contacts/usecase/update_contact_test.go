package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/cqrs/eventsourcing"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
	cmdIssuer := user.New(uuid.New())

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
					HandleCommand(ctx, gomock.Any()).
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
	cmdIssuer := user.New(uuid.New())

	t.Run("successfully update contact", func(t *testing.T) {
		uuid := uuid.New()
		cmd := CmdUpdateContact{
			ContactId: uuid.String(),
			FirstName: "John",
		}

		container.contactCmdHandler.EXPECT().
			HandleCommand(ctx, gomock.Any()).
			Times(1).
			Return(
				&domain.Contact{
					AggregateBase: eventsourcing.NewAggregateBase[domain.Contact](uuid, 2),
					FirstName:     cmd.FirstName,
				},
				nil,
			)

		updatedContact, err := contactUpdater.Update(ctx, cmd, cmdIssuer)
		assert.NoError(t, err)
		assert.Equal(t, cmd.ContactId, updatedContact.AggregateId().String())
		assert.Equal(t, cmd.FirstName, updatedContact.FirstName)
	})

	t.Run("contact not found", func(t *testing.T) {
		cmd := CmdUpdateContact{
			ContactId: uuid.NewString(),
			FirstName: "John",
		}

		container.contactCmdHandler.EXPECT().
			HandleCommand(ctx, gomock.Any()).
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
			HandleCommand(ctx, gomock.Any()).
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
		cmdWrongIssuer := user.New(uuid.New())

		container.contactCmdHandler.EXPECT().
			HandleCommand(ctx, gomock.Any()).
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
