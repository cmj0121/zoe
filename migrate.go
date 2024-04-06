package zoe

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

//go:embed assets/migrations/*.sql
var fs embed.FS

func MigrateUp(database string) int {
	source, err := iofs.New(fs, "assets/migrations")
	if err != nil {
		log.Error().Err(err).Msg("failed to create the migration source")
		return 1
	}

	switch m, err := migrate.NewWithSourceInstance("iofs", source, database); err {
	case nil:
		switch err := m.Up(); err {
		case nil:
		case migrate.ErrNoChange:
			log.Info().Msg("no change in the migration")
		default:
			log.Error().Err(err).Msg("failed to migrate up")
			return 1
		}
	default:
		log.Error().Err(err).Msg("failed to create the migration")
	}

	return 0
}
