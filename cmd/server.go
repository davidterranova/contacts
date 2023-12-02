package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/davidterranova/contacts/internal/contacts"
	"github.com/davidterranova/contacts/internal/contacts/adapters/graphql"
	lgrpc "github.com/davidterranova/contacts/internal/contacts/adapters/grpc"
	contactsHttp "github.com/davidterranova/contacts/internal/contacts/adapters/http"
	"github.com/davidterranova/contacts/internal/contacts/domain"
	contactsPorts "github.com/davidterranova/contacts/internal/contacts/ports"
	"github.com/davidterranova/contacts/pkg/user"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/davidterranova/cqrs/admin"
	adminHttp "github.com/davidterranova/cqrs/admin/adapters/http"
	"github.com/davidterranova/cqrs/eventsourcing"
	"github.com/davidterranova/cqrs/eventsourcing/eventrepository"
	"github.com/davidterranova/cqrs/eventsourcing/eventstream"
	"github.com/davidterranova/cqrs/pg"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type contactContainer struct {
	userFactory          eventsourcing.UserFactory
	eventRegistry        eventsourcing.EventRegistry[domain.Contact]
	eventRepository      eventsourcing.EventRepository
	eventSubscriber      eventsourcing.Subscriber[domain.Contact]
	eventPublisher       eventsourcing.Publisher[domain.Contact]
	eventStreamPublisher *eventsourcing.EventStreamPublisher[domain.Contact]
	eventStore           eventsourcing.EventStore[domain.Contact]
	contactFactory       eventsourcing.AggregateFactory[domain.Contact]
}

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "starts contacts server",
	PreRun: initConfig,
	Run:    runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	contactContainer, err := newContactContainer(ctx, cfg)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create contact container")
	}

	contactsApp, err := contactsApp(ctx, contactContainer, cfg)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create contacts app")
	}

	adminApp, err := adminApp(ctx, contactContainer)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create admin app")
	}

	go gqlAPIServer(ctx, contactsApp)
	go httpAPIServer(ctx, contactsApp, adminApp)
	go grpcServer(ctx, contactsApp)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		cancel()
	case <-ctx.Done():
	}
}

func httpAPIServer(ctx context.Context, contactsApp *contacts.App, adminApp *admin.App[domain.Contact]) {
	router := mux.NewRouter()
	router = contactsHttp.New(
		router,
		contactsApp,
		xhttp.GrantAnyFn(),
	)
	router = adminHttp.New[domain.Contact](
		router,
		adminApp,
	)

	xhttp.MountStatic(router, "/openapi/", "docs/openapi")
	xhttp.MountStatic(router, "/", "web")
	server := xhttp.NewServer(router, cfg.HTTP)

	err := server.Serve(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start http server")
	}
}

func gqlAPIServer(ctx context.Context, app *contacts.App) {
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: graphql.NewResolver(app)}))
	root := mux.NewRouter()
	root.Handle("/query", srv)
	root.Handle("/", playground.Handler("GraphQL playground", "/query"))
	server := xhttp.NewServer(root, cfg.GQL)

	err := server.Serve(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start graphQL server")
	}
}

func grpcServer(ctx context.Context, app *contacts.App) {
	listenTo := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	listener, err := net.Listen("tcp", listenTo)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to listen GRPC port")
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	lgrpc.RegisterContactsServer(grpcServer, lgrpc.NewHandler(app))
	log.Ctx(ctx).Info().Str("address", ":8282").Msg("starting GRPC server")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start GRPC server")
	}
}

func newContactContainer(ctx context.Context, cfg Config) (*contactContainer, error) {
	container := &contactContainer{}

	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.EventStoreDB.Name,
		ConnString: cfg.EventStoreDB.ConnString,
	})
	if err != nil {
		return nil, err
	}

	container.contactFactory = func() *domain.Contact {
		return domain.New()
	}

	container.userFactory = func() eventsourcing.User {
		return user.New(uuid.Nil)
	}
	// event registry
	container.eventRegistry = eventsourcing.NewEventRegistry[domain.Contact]()
	domain.RegisterEvents(container.eventRegistry)

	// event repository
	container.eventRepository = eventrepository.NewPGEventRepository(pg)

	// event stream
	pubSub := eventstream.NewInMemoryPubSub[domain.Contact](ctx, 100)
	container.eventSubscriber = pubSub
	container.eventPublisher = pubSub

	container.eventStreamPublisher = eventsourcing.NewEventStreamPublisher[domain.Contact](
		container.eventRepository,
		container.eventRegistry,
		domain.AggregateContact,
		container.userFactory,
		container.eventPublisher,
		100,
		true,
	)

	// event store
	withOutbox := true
	container.eventStore = eventsourcing.NewEventStore[domain.Contact](
		container.eventRepository,
		container.eventRegistry,
		container.userFactory,
		withOutbox,
	)

	// start publishing events
	go container.eventStreamPublisher.Run(ctx)

	return container, nil
}

func contactsApp(ctx context.Context, container *contactContainer, cfg Config) (*contacts.App, error) {
	contactWriteModel := eventsourcing.NewCommandHandler[domain.Contact](
		container.eventStore,
		container.contactFactory,
	)

	contactReadModel, err := contactsReadModel(
		ctx,
		cfg.ReadModelDB,
		container.eventSubscriber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create contacts read model: %w", err)
	}

	return contacts.New(contactWriteModel, contactReadModel), nil
}

func contactsReadModel(ctx context.Context, cfg pg.DBConfig, eventStreamSubscriber eventsourcing.Subscriber[domain.Contact]) (*contactsPorts.PgContactList, error) {
	// contactReadModel := ports.NewInMemoryContactList(eventStream)
	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.Name,
		ConnString: cfg.ConnString,
	})
	if err != nil {
		return nil, err
	}

	return contactsPorts.NewPgContactList(pg, eventStreamSubscriber), nil
}

func adminApp(ctx context.Context, container *contactContainer) (*admin.App[domain.Contact], error) {
	app := admin.NewApp[domain.Contact](
		container.eventRepository,
		container.eventRegistry,
		container.userFactory,
		domain.AggregateContact,
		container.contactFactory,
	)

	return app, nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
