package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// The exchange message between the honeypot service and the monitor.
type Message struct {
	Service string // the service name
	Remote  string // the remote client IP address
	Auth    *Auth  // the authentication configuration

	CreatedAt time.Time // the message created time
}

// Create a new message instance.
func New(svc string) *Message {
	return &Message{
		Service:   svc,
		CreatedAt: time.Now(),
	}
}

func MessageFromRows(rows *sql.Rows) (*Message, error) {
	var message Message
	var username sql.NullString
	var password sql.NullString
	var ns int64

	if err := rows.Scan(&message.Service, &message.Remote, &username, &password, &ns); err != nil {
		return nil, err
	}

	message.CreatedAt = time.Unix(ns/1e9, ns%1e9)

	if username.Valid && password.Valid {
		message.Auth = &Auth{
			Username: username.String,
			Password: password.String,
		}
	}

	return &message, nil
}

// Show the message as a string.
func (m Message) String() string {
	time := m.CreatedAt.UTC().Format("2006-01-02T15:04:05")
	str := fmt.Sprintf("[%v] <%s@%v> %v", time, m.Service, m.Remote, m.Auth)
	return str
}

// Show the CreatedAt time as a string.
func (m Message) CreatedTime() string {
	return m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z")
}

// Set the remote client IP address.
func (m *Message) SetRemote(remote string) *Message {
	m.Remote = remote
	return m
}

// Set the authentication configuration.
func (m *Message) SetAuth(auth *Auth) *Message {
	m.Auth = auth
	return m
}

// The customized message JSON marshaler.
func (m Message) MarshalJSON() ([]byte, error) {
	message := struct {
		Service   string // the service name
		Remote    string // the remote client IP address
		Auth      *Auth  // the authentication configuration
		CreatedAt int64  // the message created time
	}{
		Service:   m.Service,
		Remote:    m.Remote,
		Auth:      m.Auth,
		CreatedAt: m.CreatedAt.UnixNano(),
	}

	return json.Marshal(message)
}

// The interface to write the message and close the writer.
type WriteCloser interface {
	// Write the message to the writer.
	Write(msg *Message) error

	// Close the current writer.
	Close() error
}
