package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// Load the RSA key to the server configuration.
func (h *HoneypotSSH) AddRSAKey(config *ssh.ServerConfig, size int) error {
	path := filepath.Join(h.GetSSHDir(), "id_rsa")

	var signer ssh.Signer

	switch _, err := os.Stat(path); err {
	case nil:
		data, err := os.ReadFile(path)
		if err != nil {
			log.Warn().Err(err).Str("path", path).Msg("failed to read the RSA key")
			return err
		}

		if signer, err = ssh.ParsePrivateKey(data); err != nil {
			log.Warn().Err(err).Msg("failed to load the RSA signer")
			return err
		}

		log.Debug().Str("path", path).Msg("load the RSA key")
	default:
		// generate a new key
		priv, err := rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			log.Warn().Err(err).Msg("failed to generate the RSA key")
			return err
		}

		signer, err = ssh.NewSignerFromKey(priv)
		if err != nil {
			log.Warn().Err(err).Msg("failed to create the RSA signer")
			return err
		}

		// Save the private key to the file
		h.saveKey(path, priv)
	}

	config.AddHostKey(signer)
	return nil
}
