package cmd

import (
	"context"
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
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts contacts server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	eventStream := eventsourcing.NewPublisher[*domain.Contact](context.Background(), 100)
	eventStore := eventsourcing.NewEventStore[*domain.Contact]()
	contactWriteModel := eventsourcing.NewCommandHandler[*domain.Contact](
		eventStore,
		eventStream,
		func() *domain.Contact {
			return &domain.Contact{}
		},
	)

	contactReadModel := ports.NewInMemoryContactList(eventStream)
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
	router := ihttp.New(app)
	server := xhttp.NewServer(router, "", 8080)

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
	server := xhttp.NewServer(root, "", 8181)

	err := server.Serve(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start graphQL server")
	}
}

func grpcServer(ctx context.Context, app *internal.App) {
	listener, err := net.Listen("tcp", ":8282")
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

func init() {
	rootCmd.AddCommand(serverCmd)
}
