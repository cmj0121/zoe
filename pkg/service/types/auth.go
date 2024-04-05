package types

// The password-based authentication configuration for the SSH honeypot service.
type Auth struct {
	Username string // the username to authenticate
	Password string // the password to authenticate
}
