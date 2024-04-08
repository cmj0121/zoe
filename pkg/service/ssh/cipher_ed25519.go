package ssh

import (
	"crypto/ed25519"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

func (s *SSH) AddEd25519Key(conf *ssh.ServerConfig) error {
	_, priv, err := ed25519.GenerateKey(nil)

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
