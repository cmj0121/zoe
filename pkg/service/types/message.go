package types

// The exchange message between the honeypot service and the monitor.
type Message struct {
	Remote string // the remote client IP address
	Auth   *Auth  // the authentication configuration
}
