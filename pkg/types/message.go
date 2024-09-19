package types

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/database"
)

// The log that records the message of the honeypot.
type Message struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	IP      string `json:"client_ip"`
	Service string `json:"service"`

	Username *string `json:"username"`
	Password *string `json:"password"`
	Command  *string `json:"command"`
}

// Insert the message into the database.
func (m *Message) Insert() error {
	sess := database.Session()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}

	stmt := `INSERT INTO message (client_ip, service, username, password, command, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := sess.Exec(stmt, m.IP, m.Service, m.Username, m.Password, m.Command, m.CreatedAt)

	return err
}

// Iter the message from the database.
func (m *Message) Iter(ctx context.Context) <-chan *Message {
	ch := make(chan *Message, 1)

	sess := database.Session()
	base_id := math.MaxInt64
	page_size := 10

	go func() {
		defer close(ch)

		stmt := `
			SELECT id, ip, service, username, password, command, created_at
			FROM message
			WHERE id < ?
			ORDER BY id DESC
			LIMIT ?
		`

		for base_id > 0 {
			rows, err := sess.QueryContext(ctx, stmt, base_id, page_size)
			if err != nil {
				log.Warn().Err(err).Msg("failed to query the access log")
				return
			}

			for rows.Next() {
				switch msg, err := FromRow(rows); err {
				case nil:
					base_id = msg.ID
					select {
					case ch <- msg:
					case <-ctx.Done():
						return
					}
				default:
					log.Warn().Err(err).Msg("failed to parse the message")
					continue
				}
			}
		}
	}()

	return ch
}

func FromRow(rows *sql.Rows) (*Message, error) {
	var msg Message

	err := rows.Scan(&msg.ID, &msg.IP, &msg.Service, &msg.Username, &msg.Password, &msg.Command, &msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
