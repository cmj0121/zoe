package types

import (
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
		Service: svc,
		CreatedAt: time.Now(),
	}
}

// Show the message as a string.
func (m Message) String() string {
	time := m.CreatedAt.UTC().Format("2006-01-02T15:04:05")
	str := fmt.Sprintf("[%v] <%s@%v> %v", time, m.Service, m.Remote, m.Auth)
	return str
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

// The interface to write the message and close the writer.
type WriteCloser interface {
	// Write the message to the writer.
	Write(msg *Message) error

	// Close the current writer.
	Close() error
}
