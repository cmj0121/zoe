package ssh

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// Load the ED25519 key to the server configuration.
func (h *HoneypotSSH) AddED25519Key(config *ssh.ServerConfig) error {
	// generate a new key
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate the Ed25519 key")
		return err
	}

	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create the Ed25519 signer")
		return err
	}

	config.AddHostKey(signer)
	return nil
}
