package database

import (
	"database/sql"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	once      sync.Once
	defaultDB *sql.DB
)

func Init(driver, dsn string) {
	once.Do(func() {
		db, err := sql.Open(driver, dsn)
		if err != nil {
			log.Warn().Err(err).Str("driver", driver).Str("dsn", dsn).Msg("failed to open the database")
			return
		}

		defaultDB = db
		log.Info().Str("driver", driver).Str("dsn", dsn).Msg("open the database")
	})
}

func Session() *sql.DB {
	return defaultDB
}
