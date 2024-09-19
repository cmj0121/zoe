package honeypot

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/cmj0121/zoe/pkg/honeypot/ssh"
)

type HoneyPot interface {
	// Run the honeypot service that listens on the port and accepts the incoming connection.
	Run(context.Context) error
}

type Service struct {
	Name string `short:"n" help:"The name of the honeypot service" default:"ssh"`

	Service HoneyPot     `kong:"-"`
	Viper   *viper.Viper `kong:"-"`
}

func (s *Service) UnmarshalText(text []byte) error {
	s.Name = string(text)

	switch s.Name {
	case "ssh":
		s.Service = ssh.New()
	default:
		err := fmt.Errorf("unknown honeypot service: %s", s.Name)
		return err
	}
	return nil
}

func (s *Service) Run(ctx context.Context) error {
	var v *viper.Viper

	switch s.Viper {
	case nil:
		v = viper.New()
	default:
		v = s.Viper.Sub(s.Name)
	}

	// setup the default configuration
	v.SetDefault("bind", ":2022")
	v.SetDefault("server", "SSH-2.0-Open")
	v.SetDefault("max_retry", 3)
	v.SetDefault("homedir", "~")
	v.SetDefault("prompt", "$ ")
	v.SetDefault("cipher", []string{"ssh-ed25519", "rsa-sha2-256", "rsa-sha2-512"})

	if err := v.Unmarshal(s.Service); err != nil {
		log.Warn().Err(err).Msg("failed to unmarshal the service configuration")
		return err
	}

	return s.Service.Run(ctx)
}
