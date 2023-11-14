package http

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/davidterranova/contacts/internal/contacts/usecase"
	"github.com/davidterranova/contacts/pkg/auth"
	"github.com/davidterranova/contacts/pkg/user"
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
		container.handler.Use(appendUserToContextMiddleware(user.New(uuid.New())))
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
		user               user.User
		returnedAppContact *domain.Contact
		returnedAppErr     error
		expectedStatus     int
	}{
		{
			name:               "Create contact",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "jdoe@contact.local", "phone": "+15555555555"}`),
			user:               user.New(uuid.New()),
			returnedAppContact: domain.New(),
			returnedAppErr:     nil,
			expectedStatus:     http.StatusCreated,
		},
		{
			name:               "bad request",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "invalid email", "phone": "+15555555555"}`),
			user:               user.New(uuid.New()),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInvalidCommand,
			expectedStatus:     http.StatusBadRequest,
		},
		{
			name:               "internal server error",
			requestBodyContent: json.RawMessage(`{"first_name": "John", "last_name": "Doe", "email": "jdoe@contact.local", "phone": "+15555555555"}`),
			user:               user.New(uuid.New()),
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
			if c.user.IsAuthenticatedOrSystem() {
				container.handler.Use(appendUserToContextMiddleware(c.user))
			}
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
		user               user.User
		returnedAppContact *domain.Contact
		returnedAppErr     error
		expectedStatus     int
	}{
		{
			name:               "partial",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			user:               user.New(uuid.New()),
			returnedAppContact: domain.New(),
			returnedAppErr:     nil,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "not found",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			user:               user.New(uuid.New()),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrNotFound,
			expectedStatus:     http.StatusNotFound,
		},
		{
			name:               "bad request",
			requestBodyContent: json.RawMessage(`{"email": "invalid email"}`),
			user:               user.New(uuid.New()),
			returnedAppContact: nil,
			returnedAppErr:     usecase.ErrInvalidCommand,
			expectedStatus:     http.StatusBadRequest,
		},
		{
			name:               "internal server error",
			requestBodyContent: json.RawMessage(`{"first_name": "John"}`),
			user:               user.New(uuid.New()),
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
			if c.user.IsAuthenticatedOrSystem() {
				container.handler.Use(appendUserToContextMiddleware(c.user))
			}
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
		user           user.User
		returnedAppErr error
		expectedStatus int
	}{
		{
			name:           "ok",
			user:           user.New(uuid.New()),
			returnedAppErr: nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "not found",
			user:           user.New(uuid.New()),
			returnedAppErr: usecase.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal server error",
			user:           user.New(uuid.New()),
			returnedAppErr: usecase.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			container := testContainer(t)
			if c.user.IsAuthenticatedOrSystem() {
				container.handler.Use(appendUserToContextMiddleware(c.user))
			}
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
		handler: New(app, nil),
	}
}

func appendUserToContextMiddleware(u user.User) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), auth.RequestCtxUserKey, u)))
		})
	}
}
