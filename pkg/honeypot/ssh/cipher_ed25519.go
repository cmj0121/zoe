package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// Load the ED25519 key to the server configuration.
func (h *HoneypotSSH) AddED25519Key(config *ssh.ServerConfig) (err error) {
	path := filepath.Join(h.GetSSHDir(), "id_ed25519")

	var signer ssh.Signer

	switch _, err := os.Stat(path); err {
	case nil:
		data, err := os.ReadFile(path)
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("Failed to read the Ed25519 key")
			return err
		}

		if signer, err = ssh.ParsePrivateKey(data); err != nil {
			log.Warn().Err(err).Msg("Failed to load the Ed25519 signer")
			return err
		}

		log.Debug().Str("path", path).Msg("Load the Ed25519 key")
	default:
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to generate the Ed25519 key")
			return err
		}

		if signer, err = ssh.NewSignerFromKey(priv); err != nil {
			log.Warn().Err(err).Msg("Failed to create the Ed25519 signer")
			return err
		}

		// Save the private key to the file
		h.saveKey(path, priv)
	}

	config.AddHostKey(signer)
	return nil
}
