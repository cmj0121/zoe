package ssh

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

var (
	SVC_NAME = "ssh"
)

// The cipher suite for the SSH honeypot service
type Cipher struct {
	HostKeys []string // The server hostkey algorithms
}

// The honeypot service for SSH
type SSH struct {
	ch chan<- *types.Message `-` // The channel to send the message

	*types.Auth
	Cipher // the cipher suite for the SSH service

	Bind   string // the address to listen
	Banner string // the service's banner to show to the client
}

// Run the SSH honeypot service
func (s *SSH) Run(ch chan<- *types.Message) error {
	log.Info().Msg("running the SSH honeypot service ...")
	s.ch = ch

	config := &ssh.ServerConfig{
		MaxAuthTries:  3,
		ServerVersion: s.Banner,
		PasswordCallback: func(conn ssh.ConnMetadata, pwd []byte) (*ssh.Permissions, error) {
			username := conn.User()
			password := string(pwd)
			remote := conn.RemoteAddr().String()
			client_ip := strings.Split(remote, ":")[0]

			// send the authentication message
			s.sendAuthMessage(username, password, client_ip)

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
	for _, suite := range s.Cipher.HostKeys {
		switch suite {
		case "ssh-ed25519":
			if err := s.AddEd25519Key(conf); err != nil {
				allErrors = append(allErrors, err)
			}
		case "rsa-sha2-256":
			if err := s.AddRsaSha256Key(conf); err != nil {
				allErrors = append(allErrors, err)
			}
		case "rsa-sha2-512":
			if err := s.AddRsaSha512Key(conf); err != nil {
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

	remote := conn.RemoteAddr().String()
	log.Info().Str("service", SVC_NAME).Str("remote", remote).Str("bind", s.Bind).Msg("handling the connection ...")
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
		log.Info().Str("service", SVC_NAME).Interface("channel", newChannel).Msg("new channel")
		go s.handleChannel(newChannel)
	}

	log.Info().Str("service", SVC_NAME).Str("bind", s.Bind).Msg("closing the connection ...")
}

// Handle the incoming SSH channel
func (s SSH) handleChannel(newChannel ssh.NewChannel) {
	defer log.Info().Str("service", SVC_NAME).Msg("closing the channel ...")

	// Accept the channel
	channel, reqs, err := newChannel.Accept()
	if err != nil {
		log.Error().Err(err).Msg("failed to accept the channel")
		return
	}
	defer channel.Close()

	for req := range reqs {
		log.Info().Str("service", SVC_NAME).Str("type", req.Type).Msg("handling the request ...")

		switch req.Type {
		case "env":
			// ignore the environment request
			if req.WantReply {
				req.Reply(true, nil)
			}
			log.Info().Str("service", SVC_NAME).Str("type", req.Type).Msg("ignoring the request")
		case "subsystem":
			// reject the subsystem request
			if req.WantReply {
				req.Reply(false, nil)
			}
			log.Info().Str("service", SVC_NAME).Str("type", req.Type).Msg("rejecting the request")
		case "exec":
			command := string(req.Payload[4:])
			result := s.handleCommand(command)

			// reply the command result
			if req.WantReply {
				req.Reply(true, nil)
			}
			channel.Write([]byte(result))
			log.Info().Str("service", SVC_NAME).Str("type", req.Type).Str("command", command).Msg("executing the command")

			// close the channel
			return
		default:
			log.Warn().Str("service", SVC_NAME).Interface("req", req).Msg("unsupported request")
			req.Reply(false, nil)
		}
	}
}

func (s *SSH) handleCommand(command string) string {
	var result []string

	commands := strings.Split(command, ";")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		switch cmd {
		default:
			result = append(result, fmt.Sprintf("bash: %s: command not found", command))
		}
	}

	return strings.Join(result, "\n")
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
