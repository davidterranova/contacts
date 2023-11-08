package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/usecase"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"go.uber.org/mock/gomock"
)

type container struct {
	app     *MockApp
	handler *mux.Router
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("list multiple contacts", func(t *testing.T) {
		t.Parallel()

		container := testContainer(t)
		container.app.EXPECT().
			ListContacts(gomock.Any(), gomock.Any()).
			Return(
				[]*domain.Contact{
					domain.New(),
					domain.New(),
				},
				nil,
			)

		apitest.New().
			Report(apitest.SequenceDiagram()).
			Handler(container.handler).
			Get("/v1/contacts").
			Expect(t).
			Status(http.StatusOK).
			Assert(jsonpath.Len("$", 2)).
			End()
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		requestBodyContent json.RawMessage
		returnedAppContact *domain.Contact
		returnedAppErr     error
		expectedStatus     int
	}{
		{
			name:               "Create contact",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "jdoe@contact.local", "phone": "+15555555555"}`),
			returnedAppContact: domain.New(),
			returnedAppErr:     nil,
			expectedStatus:     http.StatusCreated,
		},
		{
			name:               "bad request",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "invalid email", "phone": "+15555555555"}`),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInvalidCommand,
			expectedStatus:     http.StatusBadRequest,
		},
		{
			name:               "internal server error",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "jdoe@contact.local", "phone": "+15555555555"}`),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInternal,
			expectedStatus:     http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			container := testContainer(t)
			container.app.EXPECT().
				CreateContact(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(c.returnedAppContact, c.returnedAppErr)

			apitest.New().
				Report(apitest.SequenceDiagram()).
				Handler(container.handler).
				Post("/v1/contacts").
				JSON(c.requestBodyContent).
				Expect(t).
				Status(c.expectedStatus).
				End()
		})
	}
}

func TestUpdateContact(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		requestContactId   string
		requestBodyContent json.RawMessage
		returnedAppContact *domain.Contact
		returnedAppErr     error
		expectedStatus     int
	}{
		{
			name:               "partial",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			returnedAppContact: domain.New(),
			returnedAppErr:     nil,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "not found",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrNotFound,
			expectedStatus:     http.StatusNotFound,
		},
		{
			name:               "bad request",
			requestBodyContent: json.RawMessage(`{"email": "invalid email"}`),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInvalidCommand,
			expectedStatus:     http.StatusBadRequest,
		},
		{
			name:               "internal server error",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInternal,
			expectedStatus:     http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			container := testContainer(t)
			container.app.EXPECT().
				UpdateContact(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(c.returnedAppContact, c.returnedAppErr)

			apitest.New().
				Report(apitest.SequenceDiagram()).
				Handler(container.handler).
				Putf("/v1/contacts/%s", uuid.NewString()).
				JSON(c.requestBodyContent).
				Expect(t).
				Status(c.expectedStatus).
				End()
		})
	}
}

func TestDeleteContact(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		returnedAppErr error
		expectedStatus int
	}{
		{
			name:           "ok",
			returnedAppErr: nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "not found",
			returnedAppErr: usecase.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal server error",
			returnedAppErr: usecase.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			container := testContainer(t)
			container.app.EXPECT().
				DeleteContact(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(c.returnedAppErr)

			apitest.New().
				Report(apitest.SequenceDiagram()).
				Handler(container.handler).
				Deletef("/v1/contacts/%s", uuid.NewString()).
				Expect(t).
				Status(c.expectedStatus).
				End()
		})
	}
}

func testContainer(t *testing.T) *container {
	t.Helper()

	controller := gomock.NewController(t)
	app := NewMockApp(controller)

	return &container{
		app:     app,
		handler: New(app, xhttp.GrantAnyFn()),
	}
}
