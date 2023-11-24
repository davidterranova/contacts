package cmd

import (
	"encoding/json"

	"github.com/davidterranova/contacts/pkg/xhttp"
	"github.com/davidterranova/contacts/pkg/xlogs"
	"github.com/davidterranova/cqrs/pg"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Config struct {
	Log          xlogs.LoggerConfig `envconfig:"LOG"`
	HTTP         xhttp.HTTPConfig   `envconfig:"HTTP"`
	GQL          xhttp.HTTPConfig   `envconfig:"GQL"`
	GRPC         xhttp.HTTPConfig   `envconfig:"GRPC"`
	EventStoreDB pg.DBConfig        `envconfig:"EVENT_STORE_DB"`
	ReadModelDB  pg.DBConfig        `envconfig:"READ_MODEL_DB"`
}

var cfg Config

// configCmd prints the configuration of the program
var configCmd = &cobra.Command{
	Use:    "config",
	Short:  "print the filestorage active configuration",
	PreRun: initConfig,
	Run:    runConfig,
}

func runConfig(cmd *cobra.Command, args []string) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("failed to marshal config")
	}
	log.
		Info().
		Msg(string(data))
}

func initConfig(cmd *cobra.Command, args []string) {
	// processing environment variables
	err := envconfig.Process("CONTACT", &cfg)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("failed to parse environment variables")
	}
	xlogs.Setup(
		xlogs.WithFormatter(cfg.Log.LogFormatter),
		xlogs.WithLevel(cfg.Log.LogLevel),
		xlogs.WithHostname(cfg.Log.Host),
	)
}

func init() {
	rootCmd.AddCommand(configCmd)
}
