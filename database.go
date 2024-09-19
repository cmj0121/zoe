package zoe

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/database"
)

var (
	//go:embed assets/migrations/*.sql
	migrations embed.FS
)

// The persistence layer of the honeypot holding the access logs.
type Database struct {
	Driver string `name:"driver" help:"The database driver" default:"sqlite3"`
	DSN    string `name:"dsn" help:"The data source name" default:"zoe.db"`
}

func (db *Database) Init() {
	database.Init(db.Driver, db.DSN)
}

// Migrate the database schema to the latest version.
func (db *Database) Migrate() error {
	source, err := iofs.New(migrations, "assets/migrations")
	if err != nil {
		log.Info().Msg("no migration files found")
		return err
	}

	database := fmt.Sprintf("%s://%s", db.Driver, db.DSN)
	m, err := migrate.NewWithSourceInstance("iofs", source, database)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create the migration instance")
		return err
	}

	switch err := m.Up(); err {
	case nil, migrate.ErrNoChange:
		log.Debug().Msg("migrate the database schema")
	default:
		log.Warn().Err(err).Msg("failed to migrate the database schema")
	}

	return nil
}
