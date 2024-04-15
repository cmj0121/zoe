package ssh

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var (
	SVC_NAME   = "ssh"
	SHELL_NAME = "shell"
)

// The cipher suite for the SSH honeypot service
type Cipher struct {
	HostKeys []string // The server hostkey algorithms
}

// The honeypot service for SSH
type SSH struct {
	ch             chan<- *types.Message `-` // The channel to send the message
	*term.Terminal `-`                   // The terminal for the SSH service

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
		client_ip := strings.Split(remote, ":")[0]
		go s.handleChannel(newChannel, client_ip)
	}

	log.Info().Str("service", SVC_NAME).Str("bind", s.Bind).Msg("closing the connection ...")
}

// Handle the incoming SSH channel
func (s SSH) handleChannel(newChannel ssh.NewChannel, remote string) {
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
			result := s.handleCommand(command, remote)

			// reply the command result
			if req.WantReply {
				req.Reply(true, nil)
			}
			channel.Write([]byte(result))
			log.Info().Str("service", SVC_NAME).Str("type", req.Type).Str("command", command).Msg("executing the command")

			// close the channel
			return
		case "pty-req":
			s.Terminal = term.NewTerminal(channel, s.prompt())
			req.Reply(true, nil)
		case "shell":
			go s.handleShell(channel, s.Terminal, remote)
			req.Reply(true, nil)
		default:
			log.Warn().Str("service", SVC_NAME).Interface("req", req).Msg("unsupported request")
			req.Reply(false, nil)
		}
	}
}

func (s *SSH) handleShell(channel ssh.Channel, term *term.Terminal, remote string) {
	defer log.Info().Str("service", SVC_NAME).Msg("closing the shell ...")
	defer channel.Close()

	term.SetPrompt(s.prompt())
	for {
		line, err := term.ReadLine()
		if err != nil {
			log.Warn().Err(err).Msg("failed to read the line")
			break
		}

		switch line {
		case "":
		case "exit":
			result := s.handleCommand(line, remote)
			term.Write([]byte(result))
			return
		default:
			result := s.handleCommand(line, remote)
			term.Write([]byte(result))
		}
	}
}

func (s *SSH) handleCommand(command, remote string) string {
	var result []string

	commands := strings.Split(command, ";")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)

		switch cmd {
		case "exit":
			result = append(result, "logout")
		default:
			s.sendShellMessage(cmd, remote)
			result = append(result, fmt.Sprintf("bash: %s: command not found", command))
		}
	}

	return strings.Join(result, "\n") + "\n"
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

func (s *SSH) sendShellMessage(command string, remote string) {
	message := types.New(SHELL_NAME)
	message.SetRemote(remote)
	message.Command = &command

	s.ch <- message
}

// The prompt for the SSH service
func (s *SSH) prompt() string {
	return fmt.Sprintf("%s: ~ $ ", s.Auth.Username)
}
