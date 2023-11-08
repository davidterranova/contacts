package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/user"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type container struct {
	contactCmdHandler *MockContactCmdHandler
}

func testContainer(t *testing.T) *container {
	t.Helper()

	controller := gomock.NewController(t)
	return &container{
		contactCmdHandler: NewMockContactCmdHandler(controller),
	}
}

func TestCreateContac(t *testing.T) {
	t.Parallel()

	testCreateContactValidation(t)
	testCreateContact(t)
}

func testCreateContactValidation(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactCreator := NewCreateContact(container.contactCmdHandler)
	cmdIssuer := user.New(uuid.New(), user.UserTypeAuthenticated)

	testCases := []struct {
		name          string
		command       CmdCreateContact
		expectedError error
	}{
		{
			name: "valid command",
			command: CmdCreateContact{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "test@contact.local",
				Phone:     "+33612345678",
			},
			expectedError: nil,
		},
		{
			name: "invalid command: missing email address",
			command: CmdCreateContact{
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+33612345678",
			},
			expectedError: ErrInvalidCommand,
		},
		{
			name: "invalid command: invalid phone number",
			command: CmdCreateContact{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "test@contact.local",
				Phone:     "0612345678",
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

			_, err := contactCreator.Create(ctx, tc.command, cmdIssuer)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testCreateContact(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactCreator := NewCreateContact(container.contactCmdHandler)
	cmdIssuer := user.New(uuid.New(), user.UserTypeAuthenticated)

	t.Run("successful contact creation", func(t *testing.T) {
		cmd := CmdCreateContact{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@contact.local",
			Phone:     "+33612345678",
		}

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				&domain.Contact{
					FirstName: cmd.FirstName,
					LastName:  cmd.LastName,
					Email:     cmd.Email,
					Phone:     cmd.Phone,
					CreatedBy: cmdIssuer,
				},
				nil,
			)

		createdContact, err := contactCreator.Create(ctx, cmd, cmdIssuer)
		assert.NoError(t, err)
		assert.Equal(t, cmd.FirstName, createdContact.FirstName)
		assert.Equal(t, cmd.LastName, createdContact.LastName)
		assert.Equal(t, cmd.Email, createdContact.Email)
		assert.Equal(t, cmd.Phone, createdContact.Phone)
		assert.Equal(t, cmdIssuer, createdContact.CreatedBy)
	})

	t.Run("repository unexpected error", func(t *testing.T) {
		cmd := CmdCreateContact{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@contact.local",
			Phone:     "+33612345678",
		}

		container.contactCmdHandler.EXPECT().
			Handle(gomock.Any()).
			Times(1).
			Return(
				nil,
				errors.New("unexpected error"),
			)

		contact, err := contactCreator.Create(ctx, cmd, cmdIssuer)
		assert.ErrorIs(t, err, ErrInternal)
		assert.Nil(t, contact)
	})
}
