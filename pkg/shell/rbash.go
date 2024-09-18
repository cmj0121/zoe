// The mock restrict-bash (rbash) service.
package shell

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// The restricted bash shell that provides the limited bash shell.
// It is the semi-interactive shell that accepts the command and returns the output.
type RBash struct{}

// New creates a new RBash instance that provides the restricted bash shell.
func New() *RBash {
	return &RBash{}
}

// Exec the command and return the output as the rbash shell.
func (r *RBash) Exec(command string) string {
	cmds := strings.Split(command, "; ")

	result := []string{}
	for _, cmd := range cmds {
		// disallow the I/O redirection
		switch {
		case strings.Contains(cmd, ">"):
			result = append(result, "bash: I/O redirection is not allowed")
			continue
		case strings.Contains(cmd, "<"):
			result = append(result, "bash: I/O redirection is not allowed")
			continue
		case strings.Contains(cmd, "|"):
			result = append(result, "bash: I/O redirection is not allowed")
			continue
		}

		if cmd = strings.TrimSpace(cmd); cmd == "" {
			log.Debug().Msg("skip the empty command")
			continue
		}

		args := strings.Split(cmd, " ")
		result = append(result, r.exec(args[0], args[1:]...))
	}

	return strings.Join(result, "\n")
}

// Execute the command and return the output as the restricted bash shell.
func (r *RBash) exec(command string, args ...string) string {
	var output string
	log.Info().Str("command", command).Strs("args", args).Msg("exec the command")

	switch command {
	case "ls":
		output = ".ssh"
	case "pwd":
		output = "/home/nobody"
	case "whoami":
		output = "nobody"
	case "echo":
		output = strings.Join(args, " ")
	default:
		output = fmt.Sprintf("bash: %s: command not found", command)
	}

	return output
}
