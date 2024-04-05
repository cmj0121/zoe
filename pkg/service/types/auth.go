package types

import (
	"fmt"
)

// The password-based authentication configuration for the SSH honeypot service.
type Auth struct {
	Username string // the username to authenticate
	Password string // the password to authenticate
}

func (a Auth) String() string {
	str := fmt.Sprintf("%s:%s", a.Username, a.Password)
	return str
}
