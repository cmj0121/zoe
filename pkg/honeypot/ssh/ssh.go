package ssh

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// The SSH-based honeypot service that provides the semi-interactive shell.
type HoneypotSSH struct {
	Bind     string `short:"b" help:"The address to bind" default:":2022"`
	Server   string `help:"The server version" default:"SSH-2.0-Open"`
	MaxRetry int    `short:"r" help:"The maximum number of the retry" default:"3"`

	Username *string  `short:"u" help:"The authorized username"`
	Password *string  `short:"P" help:"The authorized password"`
	Cipher   []string `help:"The list of the cipher" default:"ssh-ed25519,rsa-sha2-256,rsa-sha2-512"`
}

func New() *HoneypotSSH {
	return &HoneypotSSH{}
}

// Run the honeypot service that listens on the port and accepts the incoming SSH connection.
func (h *HoneypotSSH) Run(ctx context.Context) error {
	config := &ssh.ServerConfig{
		MaxAuthTries:  h.MaxRetry,
		ServerVersion: h.Server,
		PasswordCallback: func(conn ssh.ConnMetadata, bytes []byte) (*ssh.Permissions, error) {
			username := conn.User()
			password := string(bytes)

			switch {
			case h.Username == nil:
				log.Debug().Msg("no authorized username, always reject the connection")
				return nil, fmt.Errorf("no authorized username")
			case username == *h.Username && h.Password == nil:
				log.Debug().Msg("no authorized password, always accept the connection")
			case username != *h.Username || password != *h.Password:
				log.Debug().Msg("invalid username or password")
				return nil, fmt.Errorf("invalid username or password")
			}

			log.Info().Str("username", username).Str("password", password).Msg("accept the SSH connection")
			return nil, nil
		},
	}

	h.AddHostKey(config)
	return h.run(ctx, config)
}

// Start the SSH service with the given configuration.
func (h *HoneypotSSH) run(ctx context.Context, cfg *ssh.ServerConfig) error {
	listener, err := net.Listen("tcp", h.Bind)
	if err != nil {
		log.Warn().Err(err).Msg("failed to listen on the address")
		return err
	}
	// pass the listener and generate the connection handler
	handler := h.handleTCPConn(ctx, listener)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("the service is shutting down")
			return nil
		case conn := <-handler:
			go h.handleSSHConn(ctx, conn, cfg)
		}
	}
}

// The new connection handler that accepts the incoming TCP connection.
func (h *HoneypotSSH) handleTCPConn(ctx context.Context, listener net.Listener) <-chan net.Conn {
	ch := make(chan net.Conn)

	log.Info().Str("bind", h.Bind).Msg("the SSH service is listening on the address")
	go func() {
		defer listener.Close()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("the service is shutting down")
				break
			default:
				conn, err := listener.Accept()
				if err != nil {
					log.Warn().Err(err).Msg("failed to accept the incoming connection")
					continue
				}

				ch <- conn
			}
		}
	}()

	return ch
}

// Handle the SSH connection with the given configuration.
func (h *HoneypotSSH) handleSSHConn(ctx context.Context, conn net.Conn, cfg *ssh.ServerConfig) {
	defer conn.Close()

	remote := conn.RemoteAddr().String()
	log.Info().Str("remote", remote).Str("bind", h.Bind).Msg("accepted the incoming TCP connection")

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, cfg)
	if err != nil {
		log.Warn().Err(err).Msg("failed to handshake the SSH connection")
		return
	}

	client := sshConn.RemoteAddr().String()
	log.Info().Str("client", client).Str("remote", remote).Msg("accepted the incoming SSH connection")
	// discard the requests
	go ssh.DiscardRequests(reqs)

	for channel := range chans {
		go h.handleSSHChannel(ctx, channel)
	}
}

// Handle the SSH channel with the given configuration.
func (h *HoneypotSSH) handleSSHChannel(ctx context.Context, channel ssh.NewChannel) {
	ch, reqs, err := channel.Accept()
	if err != nil {
		log.Warn().Err(err).Msg("failed to accept the SSH channel")
		return
	}

	defer ch.Close()

	for req := range reqs {
		switch req.Type {
		default:
			log.Warn().Str("type", req.Type).Msg("unsupported request")

			if err := req.Reply(false, nil); err != nil {
				log.Warn().Err(err).Msg("failed to reply the request")
				continue
			}
		}
	}
}
