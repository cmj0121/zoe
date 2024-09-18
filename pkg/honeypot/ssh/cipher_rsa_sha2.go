package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// Load the RSA key to the server configuration.
func (h *HoneypotSSH) AddRSAKey(config *ssh.ServerConfig, size int) error {
	// generate a new key
	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate the RSA key")
		return err
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create the RSA signer")
		return err
	}

	config.AddHostKey(signer)
	return nil
}
