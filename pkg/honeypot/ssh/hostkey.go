package ssh

import (
	"crypto"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// add the host key to the server configuration, may load from the file or generate a new one
// by default, the server will generate a new host key if the host key is not set.
func (h *HoneypotSSH) AddHostKey(config *ssh.ServerConfig) {
	for _, suite := range h.Cipher {
		switch suite {
		case ssh.KeyAlgoED25519:
			if err := h.AddED25519Key(config); err != nil {
				log.Warn().Err(err).Msg("Failed to add the Ed25519 key")
			}
		case ssh.KeyAlgoRSASHA256:
			if err := h.AddRSAKey(config, 2048); err != nil {
				log.Warn().Err(err).Msg("Failed to add the RSA key")
			}
		case ssh.KeyAlgoRSASHA512:
			if err := h.AddRSAKey(config, 4096); err != nil {
				log.Warn().Err(err).Msg("Failed to add the RSA key")
			}
		default:
			log.Warn().Str("suite", suite).Msg("unknown cipher suite")
			continue
		}
	}
}

// Get the .ssh folder of the user.
func (h *HoneypotSSH) GetSSHDir() string {
	var basedir string

	switch h.Homedir {
	case "~":
		switch dir, err := os.UserHomeDir(); err {
		case nil:
			basedir = dir
		default:
			log.Warn().Err(err).Msg("failed to get the home directory")
			basedir = "."
		}
	default:
		basedir = h.Homedir
	}

	return filepath.Join(basedir, ".ssh")
}

func (h *HoneypotSSH) saveKey(path string, priv crypto.PrivateKey) {
	switch pkey, err := ssh.MarshalPrivateKey(priv, ""); err {
	case nil:
		text := base64.StdEncoding.EncodeToString(pkey.Bytes)
		// always make the text wrapped with the maximum 70 characters
		// as the OpenSSH does
		keys := []string{"-----BEGIN OPENSSH PRIVATE KEY-----"}
		for len(text) > 0 {
			switch {
			case len(text) >= 70:
				keys = append(keys, text[:70])
				text = text[70:]
			default:
				keys = append(keys, text)
				text = ""
			}
		}
		keys = append(keys, "-----END OPENSSH PRIVATE KEY-----")
		pkey := strings.Join(keys, "\n")

		if err := os.WriteFile(path, []byte(pkey), 0600); err != nil {
			log.Warn().Err(err).Str("path", path).Msg("failed to save the Ed25519 private key")
			return
		}
	default:
		log.Warn().Err(err).Msg("failed to marshal the private key")
		return
	}
}
