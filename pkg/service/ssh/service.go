package ssh

import (
	"crypto/ed25519"
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
	*types.Auth

	Bind    string   // the address to listen
	Banner  string   // the service's banner to show to the client
	Ciphers []string // the list of ciphers to use

	signer ssh.Signer // The private key to use for the host key
}

// Run the SSH honeypot service
func (s *SSH) Run(ch <-chan types.Message) error {
	log.Info().Msg("running the SSH honeypot service ...")

	// Set-up the SSH server configuration
	signer, err := s.getSigner()
	if err != nil {
		log.Error().Err(err).Msg("failed to get the signer")
		return err
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, pwd []byte) (*ssh.Permissions, error) {
			username := conn.User()
			password := string(pwd)

			if s.Auth == nil {
				log.Debug().Str("service", SVC_NAME).Msg("disabled the authentication")
				return nil, fmt.Errorf("authentication disabled")
			}

			client := conn.RemoteAddr().String()
			if username == s.Auth.Username && password == s.Auth.Password {
				log.Info().Str("service", SVC_NAME).Str("client", client).Msg("authenticated")
				return nil, nil
			}

			log.Info().Str("service", SVC_NAME).Str("client", client).Msg("authentication failed")
			return nil, fmt.Errorf("authentication failed: %v", conn.User())
		},
	}
	config.AddHostKey(signer)

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

// Get the SSH key signature
func (s *SSH) getSigner() (ssh.Signer, error) {
	if s.signer == nil {
		_, priv, err := ed25519.GenerateKey(nil)

		if err != nil {
			log.Error().Err(err).Msg("failed to generate the private key")
			return nil, err
		}

		signer, err := ssh.NewSignerFromKey(priv)
		if err != nil {
			log.Error().Err(err).Msg("failed to create the signer")
			return nil, err
		}

		s.signer = signer
	}

	return s.signer, nil
}
