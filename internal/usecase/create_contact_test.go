package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type container struct {
	contactRepo *MockContactRepository
}

func testContainer(t *testing.T) *container {
	t.Helper()

	controller := gomock.NewController(t)
	return &container{
		contactRepo: NewMockContactRepository(controller),
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
	contactCreator := NewCreateContact(container.contactRepo)

	testCases := []struct {
		name          string
		command       CmdCreateContact
		expectedError error
	}{
		{
			name: "valid command",
			command: CmdCreateContact{
				CreatedBy: user.New(uuid.New(), user.UserTypeAuthenticated),
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
				CreatedBy: user.New(uuid.New(), user.UserTypeAuthenticated),
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+33612345678",
			},
			expectedError: ErrInvalidCommand,
		},
		{
			name: "invalid command: invalid phone number",
			command: CmdCreateContact{
				CreatedBy: user.New(uuid.New(), user.UserTypeAuthenticated),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "test@contact.local",
				Phone:     "0612345678",
			},
			expectedError: ErrInvalidCommand,
		},
		{
			name: "invalid command: invalid uuid",
			command: CmdCreateContact{
				CreatedBy: user.NewUnauthenticated(),
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
				container.contactRepo.EXPECT().
					Create(ctx, gomock.Any()).
					Times(1).
					Return(nil, nil)
			}

			_, err := contactCreator.Create(ctx, tc.command)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func testCreateContact(t *testing.T) {
	ctx := context.Background()
	container := testContainer(t)
	contactCreator := NewCreateContact(container.contactRepo)

	t.Run("successful contact creation", func(t *testing.T) {
		cmd := CmdCreateContact{
			CreatedBy: user.New(uuid.New(), user.UserTypeAuthenticated),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@contact.local",
			Phone:     "+33612345678",
		}

		container.contactRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Times(1).
			Return(
				&domain.Contact{
					FirstName: cmd.FirstName,
					LastName:  cmd.LastName,
					Email:     cmd.Email,
					Phone:     cmd.Phone,
				},
				nil,
			)

		createdContact, err := contactCreator.Create(ctx, cmd)
		assert.NoError(t, err)
		assert.Equal(t, cmd.FirstName, createdContact.FirstName)
		assert.Equal(t, cmd.LastName, createdContact.LastName)
		assert.Equal(t, cmd.Email, createdContact.Email)
		assert.Equal(t, cmd.Phone, createdContact.Phone)
	})

	t.Run("repository unexpected error", func(t *testing.T) {
		cmd := CmdCreateContact{
			CreatedBy: user.New(uuid.New(), user.UserTypeAuthenticated),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "jdoe@contact.local",
			Phone:     "+33612345678",
		}

		container.contactRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Times(1).
			Return(
				nil,
				errors.New("unexpected error"),
			)

		contact, err := contactCreator.Create(ctx, cmd)
		assert.ErrorIs(t, err, ErrInternal)
		assert.Nil(t, contact)
	})
}
