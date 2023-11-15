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
	"github.com/davidterranova/contacts/internal/admin"
	adminPorts "github.com/davidterranova/contacts/internal/admin/ports"
	"github.com/davidterranova/contacts/internal/contacts"
	"github.com/davidterranova/contacts/internal/contacts/adapters/graphql"
	lgrpc "github.com/davidterranova/contacts/internal/contacts/adapters/grpc"
	"github.com/davidterranova/contacts/internal/contacts/domain"
	contactsPorts "github.com/davidterranova/contacts/internal/contacts/ports"

	adminHttp "github.com/davidterranova/contacts/internal/admin/adapters/http"
	contactsHttp "github.com/davidterranova/contacts/internal/contacts/adapters/http"
	"github.com/davidterranova/contacts/pkg/eventsourcing"
	"github.com/davidterranova/contacts/pkg/pg"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "starts contacts server",
	PreRun: initConfig,
	Run:    runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	contactsApp, err := contactsApp(ctx, cfg)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create contacts app")
	}

	adminApp, err := adminApp(ctx, cfg)
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

func httpAPIServer(ctx context.Context, contactsApp *contacts.App, adminApp *admin.App) {
	router := mux.NewRouter()
	router = contactsHttp.New(
		router,
		contactsApp,
		xhttp.GrantAnyFn(),
	)
	router = adminHttp.New(
		router,
		adminApp,
		nil,
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

func contactsApp(ctx context.Context, cfg Config) (*contacts.App, error) {
	eventRegistry := eventsourcing.NewRegistry[domain.Contact]()
	domain.RegisterEvents(eventRegistry)

	eventStream := eventsourcing.NewInMemoryPublisher[domain.Contact](context.Background(), 100)
	contactWriteModel, eventStreamPublisher, err := contactsWriteModel(ctx, cfg.EventStoreDB, eventRegistry, eventStream)
	if err != nil {
		return nil, fmt.Errorf("failed to create write model: %w", err)
	}

	// start publishing events
	go eventStreamPublisher.Run(ctx)

	contactReadModel, err := contactsReadModel(ctx, cfg.ReadModelDB, eventStream)
	if err != nil {
		return nil, fmt.Errorf("failed to create contacts read model: %w", err)
	}

	return contacts.New(contactWriteModel, contactReadModel), nil
}

func contactsWriteModel(ctx context.Context, cfg pg.DBConfig, eventRegistry *eventsourcing.EventRegistry[domain.Contact], eventStream eventsourcing.EventStream[domain.Contact]) (eventsourcing.CommandHandler[domain.Contact], *eventsourcing.EventStreamPublisher[domain.Contact], error) {
	// eventStore := eventsourcing.NewInMemoryEventStore[domain.Contact]()
	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.Name,
		ConnString: cfg.ConnString,
	})
	if err != nil {
		return nil, nil, err
	}

	eventStore := eventsourcing.NewPGEventStore[domain.Contact](pg, eventRegistry)
	eventPublisher := eventsourcing.NewEventStreamPublisher[domain.Contact](eventStore, eventStream, 10)

	contactWriteModel := eventsourcing.NewCommandHandler[domain.Contact](
		eventStore,
		func() *domain.Contact {
			return domain.New()
		},
	)

	return contactWriteModel, eventPublisher, nil
}

func contactsReadModel(ctx context.Context, cfg pg.DBConfig, eventStream eventsourcing.EventStream[domain.Contact]) (*contactsPorts.PgContactList, error) {
	// contactReadModel := ports.NewInMemoryContactList(eventStream)
	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.Name,
		ConnString: cfg.ConnString,
	})
	if err != nil {
		return nil, err
	}

	return contactsPorts.NewPgContactList(pg, eventStream), nil
}

func adminApp(ctx context.Context, cfg Config) (*admin.App, error) {
	adminReadModel, err := adminReadModel(ctx, cfg.EventStoreDB)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create admin read model")
	}

	return admin.New(adminReadModel), nil
}

func adminReadModel(ctx context.Context, cfg pg.DBConfig) (*adminPorts.PgListEvent, error) {
	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.Name,
		ConnString: cfg.ConnString,
	})
	if err != nil {
		return nil, err
	}

	return adminPorts.NewPgListEvent(pg), nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
