package cmd

import (
	"errors"

	"github.com/davidterranova/contacts/pkg/pg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	dbmigrate = &cobra.Command{
		Use:    "dbmigrate",
		Short:  "migrate db schemas",
		PreRun: initConfig,
	}

	dbmigrateUp = &cobra.Command{
		Use:    "up",
		Short:  "migrate db schemas up",
		Run:    runDBMigrateUp,
		PreRun: initDBMigrate,
	}

	dbmigrateDown = &cobra.Command{
		Use:    "down",
		Short:  "migrate db schemas down",
		Run:    runDBMigratDown,
		PreRun: initDBMigrate,
	}

	migrator *migrate.Migrate
	target   *string
)

func runDBMigrateUp(cmd *cobra.Command, args []string) {
	log.Info().Msg("migrating db schemas up")

	err := migrator.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("no changes to migrate")
			return
		}
		log.Fatal().Err(err).Msg("failed to migrate up")
	}
}

func runDBMigratDown(cmd *cobra.Command, args []string) {
	log.Info().Msg("migrating db schemas down")

	err := migrator.Down()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate down")
	}
}

func initDBMigrate(cmd *cobra.Command, args []string) {
	initConfig(cmd, args)

	fs, err := pg.NewMigratorFS()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migrator filesystem")
	}

	var dbConfig pg.DBConfig
	switch *target {
	case "eventstore":
		dbConfig = cfg.EventStoreDB
	default:
		log.Fatal().Str("available", "[eventstore]").Str("target", *target).Msg("unknown target")
	}

	migrator, err = migrate.NewWithSourceInstance("iofs", fs, string(dbConfig.ConnString))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migrate instance")
	}
}

func init() {
	target = dbmigrate.PersistentFlags().StringP("target", "t", "eventstore", "target migration database")

	dbmigrate.AddCommand(dbmigrateUp)
	dbmigrate.AddCommand(dbmigrateDown)
	rootCmd.AddCommand(dbmigrate)
}
