package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

func (s *SSH) AddRsaSha256Key(conf *ssh.ServerConfig) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Error().Err(err).Msg("failed to generate the private key")
		return err
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Error().Err(err).Msg("failed to create the signer")
		return err
	}

	conf.AddHostKey(signer)
	return nil
}

func (s *SSH) AddRsaSha512Key(conf *ssh.ServerConfig) error {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		log.Error().Err(err).Msg("failed to generate the private key")
		return err
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Error().Err(err).Msg("failed to create the signer")
		return err
	}

	conf.AddHostKey(signer)
	return nil
}
