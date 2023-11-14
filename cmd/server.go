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
	"github.com/davidterranova/contacts/internal"
	"github.com/davidterranova/contacts/internal/adapters/graphql"
	lgrpc "github.com/davidterranova/contacts/internal/adapters/grpc"
	"github.com/davidterranova/contacts/internal/domain"
	"github.com/davidterranova/contacts/internal/ports"

	ihttp "github.com/davidterranova/contacts/internal/adapters/http"
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

	eventRegistry := eventsourcing.NewRegistry[domain.Contact]()
	domain.RegisterEvents(eventRegistry)

	eventStream := eventsourcing.NewInMemoryPublisher[domain.Contact](context.Background(), 100)
	contactWriteModel, eventStreamPublisher, err := writeModel(ctx, cfg.EventStoreDB, eventRegistry, eventStream)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create write model")
	}

	// start publishing events
	go eventStreamPublisher.Run(ctx)

	contactReadModel, err := readModel(ctx, cfg.ReadModelDB, eventStream)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to create read model")
	}

	app := internal.New(contactWriteModel, contactReadModel)

	go gqlAPIServer(ctx, app)
	go httpAPIServer(ctx, app)
	go grpcServer(ctx, app)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		cancel()
	case <-ctx.Done():
	}
}

func httpAPIServer(ctx context.Context, app *internal.App) {
	router := ihttp.New(
		app,
		xhttp.GrantAnyFn(),
	)
	server := xhttp.NewServer(router, cfg.HTTP)

	err := server.Serve(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start http server")
	}
}

func gqlAPIServer(ctx context.Context, app *internal.App) {
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

func grpcServer(ctx context.Context, app *internal.App) {
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

func writeModel(ctx context.Context, cfg pg.DBConfig, eventRegistry *eventsourcing.Registry[domain.Contact], eventStream eventsourcing.EventStream[domain.Contact]) (eventsourcing.CommandHandler[domain.Contact], *eventsourcing.EventStreamPublisher[domain.Contact], error) {
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

func readModel(ctx context.Context, cfg pg.DBConfig, eventStream eventsourcing.EventStream[domain.Contact]) (*ports.PgContactList, error) {
	// contactReadModel := ports.NewInMemoryContactList(eventStream)

	pg, err := pg.Open(pg.DBConfig{
		Name:       cfg.Name,
		ConnString: cfg.ConnString,
	})
	if err != nil {
		return nil, err
	}

	return ports.NewPgContactList(pg, eventStream), nil
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
