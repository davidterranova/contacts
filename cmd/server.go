package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/davidterranova/contacts/internal"
	"github.com/davidterranova/contacts/internal/adapters/http"
	"github.com/davidterranova/contacts/internal/ports"
	"github.com/davidterranova/contacts/pkg/xhttp"
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

	router := http.New(app)
	server := xhttp.NewServer(router, "", 8080)

	err := server.Serve(ctx)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to start server")
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
