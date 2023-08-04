package usecase

// func TestUpdateContac(t *testing.T) {
// 	t.Parallel()

// 	testUpdateContactValidation(t)
// 	testUpdateContact(t)
// }

// func testUpdateContactValidation(t *testing.T) {
// 	ctx := context.Background()
// 	container := testContainer(t)
// 	contactUpdater := NewUpdateContact(container.contactRepo)

// 	testCases := []struct {
// 		name          string
// 		command       CmdUpdateContact
// 		expectedError error
// 	}{
// 		{
// 			name: "valid command",
// 			command: CmdUpdateContact{
// 				ContactId: uuid.NewString(),
// 				FirstName: "John",
// 				LastName:  "Doe",
// 				Email:     "test@contact.local",
// 				Phone:     "+33612345678",
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "invalid command: invalid contact id",
// 			command: CmdUpdateContact{
// 				ContactId: "invalid-uuid",
// 				FirstName: "John",
// 			},
// 			expectedError: ErrInvalidCommand,
// 		},
// 		{
// 			name: "invalid command: missing contact id",
// 			command: CmdUpdateContact{
// 				FirstName: "John",
// 			},
// 			expectedError: ErrInvalidCommand,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			if tc.expectedError == nil {
// 				container.contactRepo.EXPECT().
// 					Update(ctx, gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(nil, nil)
// 			}

// 			_, err := contactUpdater.Update(ctx, tc.command)
// 			assert.ErrorIs(t, err, tc.expectedError)
// 		})
// 	}
// }

// func testUpdateContact(t *testing.T) {
// 	ctx := context.Background()
// 	container := testContainer(t)
// 	contactUpdater := NewUpdateContact(container.contactRepo)

// 	t.Run("successfully update contact", func(t *testing.T) {
// 		uuid := uuid.New()
// 		cmd := CmdUpdateContact{
// 			ContactId: uuid.String(),
// 			FirstName: "John",
// 		}

// 		container.contactRepo.EXPECT().
// 			Update(ctx, gomock.Any(), gomock.Any()).
// 			Times(1).
// 			Return(
// 				&domain.Contact{
// 					Id:        uuid,
// 					FirstName: cmd.FirstName,
// 				},
// 				nil,
// 			)

// 		updatedContact, err := contactUpdater.Update(ctx, cmd)
// 		assert.NoError(t, err)
// 		assert.Equal(t, cmd.ContactId, updatedContact.Id.String())
// 		assert.Equal(t, cmd.FirstName, updatedContact.FirstName)
// 	})

// 	t.Run("contact not found", func(t *testing.T) {
// 		cmd := CmdUpdateContact{
// 			ContactId: uuid.NewString(),
// 			FirstName: "John",
// 		}

// 		container.contactRepo.EXPECT().
// 			Update(ctx, gomock.Any(), gomock.Any()).
// 			Times(1).
// 			Return(
// 				nil,
// 				ports.ErrNotFound,
// 			)

// 		updatedContact, err := contactUpdater.Update(ctx, cmd)
// 		assert.ErrorIs(t, err, ErrNotFound)
// 		assert.Nil(t, updatedContact)
// 	})

// 	t.Run("unexpected repository error", func(t *testing.T) {
// 		cmd := CmdUpdateContact{
// 			ContactId: uuid.NewString(),
// 			FirstName: "John",
// 		}

// 		container.contactRepo.EXPECT().
// 			Update(ctx, gomock.Any(), gomock.Any()).
// 			Times(1).
// 			Return(
// 				nil,
// 				errors.New("internal error"),
// 			)

// 		updatedContact, err := contactUpdater.Update(ctx, cmd)
// 		assert.ErrorIs(t, err, ErrInternal)
// 		assert.Nil(t, updatedContact)
// 	})
// }
