package ssh

import (
	"errors"
	"fmt"
	"net"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

var (
	SVC_NAME = "ssh"
)

// The honeypot service for SSH
type SSH struct {
	ch chan<- *types.Message `-` // The channel to send the message

	*types.Auth

	Bind    string   // the address to listen
	Banner  string   // the service's banner to show to the client
	Ciphers []string // the list of ciphers to use
}

// Run the SSH honeypot service
func (s *SSH) Run(ch chan<- *types.Message) error {
	log.Info().Msg("running the SSH honeypot service ...")
	s.ch = ch

	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, pwd []byte) (*ssh.Permissions, error) {
			username := conn.User()
			password := string(pwd)
			remote := conn.RemoteAddr().String()

			// send the authentication message
			s.sendAuthMessage(username, password, remote)

			if s.Auth == nil {
				log.Debug().Str("service", SVC_NAME).Msg("disabled the authentication")
				return nil, fmt.Errorf("authentication disabled")
			}

			if username == s.Auth.Username && password == s.Auth.Password {
				log.Info().Str("service", SVC_NAME).Str("client", remote).Msg("authenticated")
				return nil, nil
			}

			log.Info().Str("service", SVC_NAME).Str("client", remote).Msg("authentication failed")
			return nil, fmt.Errorf("authentication failed: %v", conn.User())
		},
	}
	if err := s.AddHostKey(config); err != nil {
		log.Error().Err(err).Msg("failed to add the host key")
		return err
	}

	// create a TCP listener
	listener, err := net.Listen("tcp", s.Bind)
	if err != nil {
		log.Error().Err(err).Str("bind", s.Bind).Msg("failed to create a listener")
		return err
	}
	defer listener.Close()

	// accept the incoming connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept the connection")
			continue
		}

		// handle the connection in a goroutine
		go s.handleConn(conn, config)
	}

	return nil
}

// Add the host key to the SSH server configuration
func (s *SSH) AddHostKey(conf *ssh.ServerConfig) error {
	allErrors := []error{}
	for _, suite := range s.Ciphers {
		switch suite {
		case "ssh-ed25519":
			if err := s.AddEd25519Key(conf); err != nil {
				allErrors = append(allErrors, err)
			}
		default:
			err := fmt.Errorf("unknown cipher suite: %s", suite)
			allErrors = append(allErrors, err)
		}
	}

	return errors.Join(allErrors...)
}

// Handle the incoming TCP connection for the SSH honeypot service
func (s *SSH) handleConn(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	log.Info().Str("service", SVC_NAME).Str("bind", s.Bind).Msg("handling the connection ...")
	// Before use, a handshake must be performed on the incoming net.Conn.
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Error().Err(err).Msg("failed to handshake")
		return
	}

	log.Info().Str("service", SVC_NAME).Str("client", sshConn.RemoteAddr().String()).Msg("handshake success")
	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// handle the channel in another goroutine
		go s.handleChannel(newChannel)
	}

	log.Info().Str("service", SVC_NAME).Str("bind", s.Bind).Msg("closing the connection ...")
}

// Handle the incoming SSH channel
func (s SSH) handleChannel(newChannel ssh.NewChannel) {
	// Accept the channel
	channel, _, err := newChannel.Accept()
	if err != nil {
		log.Error().Err(err).Msg("failed to accept the channel")
		return
	}
	defer channel.Close()
}

func (s *SSH) sendAuthMessage(username string, password string, remote string) {
	message := types.New(SVC_NAME)
	message.SetRemote(remote)
	message.SetAuth(&types.Auth{
		Username: username,
		Password: password,
	})

	s.ch <- message
}
