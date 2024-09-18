package ssh

import (
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
