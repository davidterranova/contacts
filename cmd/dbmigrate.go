package cmd

import (
	"embed"
	"errors"

	lpg "github.com/davidterranova/contacts/pkg/pg"
	"github.com/davidterranova/cqrs/pg"
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

	var (
		dbConfig pg.DBConfig
		fs       embed.FS
		path     string
		err      error
	)

	migrations := pg.NewMigrations()
	migrations.Append("eventstore", "migrations", pg.EventSourcingFS)
	migrations.Append("readmodel", "migrations", lpg.ReadModelFS)

	switch *target {
	case "eventstore":
		log.Info().Str("target", "eventstore").Msg("migration configuration loaded")
		dbConfig = cfg.EventStoreDB
		fs, path, err = migrations.Get("eventstore")
	case "readmodel":
		log.Info().Str("target", "readmodel").Msg("migration configuration loaded")
		dbConfig = cfg.ReadModelDB
		fs, path, err = migrations.Get("readmodel")
	default:
		log.Fatal().Str("available", "[eventstore]").Str("target", *target).Msg("unknown target")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get migration configuration")
	}

	driver, err := pg.NewMigratorFS(fs, path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migrator filesystem")
	}

	migrator, err = migrate.NewWithSourceInstance("iofs", driver, string(dbConfig.ConnString))
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
