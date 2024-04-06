package config

import (
	"database/sql"

	"github.com/cmj0121/zoe/pkg/service/types"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Save the logger and the writer to the console
type SQLiteLogger struct {
	*sql.DB
	path string
}

func NewSQLiteLogger(path string) (*SQLiteLogger, error) {
	switch db, err := sql.Open("sqlite3", path); err {
	case nil:
		logger := &SQLiteLogger{
			DB:   db,
			path: path,
		}
		return logger, nil
	default:
		return nil, err
	}
}

func (s *SQLiteLogger) Write(msg *types.Message) error {
	stmt, err := s.Prepare("INSERT OR IGNORE INTO message (service, client_ip, username, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare the statement")
		return err
	}

	switch msg.Auth {
	case nil:
		null := sql.NullString{}
		_, err = stmt.Exec(msg.Service, null, null)
	default:
		_, err = stmt.Exec(msg.Service, msg.Remote, msg.Auth.Username, msg.Auth.Password)
	}

	return err
}

func (s *SQLiteLogger) Close() error {
	return s.Close()
}
