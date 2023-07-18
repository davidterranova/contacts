package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/davidterranova/contacts/internal"
	"github.com/davidterranova/contacts/internal/adapters/graphql"
	ihttp "github.com/davidterranova/contacts/internal/adapters/http"
	"github.com/davidterranova/contacts/internal/ports"
	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts contacts server",
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := internal.New(ports.NewInMemoryContactRepository())

	go gqlAPIServer(ctx, app)
	go httpAPIServer(ctx, app)

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
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start server")
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
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start server")
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
