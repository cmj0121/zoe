package zoe

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func MigrateUp(database, folder string) int {
	log.Info().Str("database", database).Str("folder", folder).Msg("execute the migration")

	switch m, err := migrate.New(fmt.Sprintf("file://%s", folder), database); err {
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
