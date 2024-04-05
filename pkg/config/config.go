package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/cmj0121/zoe/pkg/service"
)

// The configuration of the Zoe service including
//   - how to log the behavior
//   - the web interface to show the behavior
//   - the honeypot service to catch the behavior
type Config struct {
	Loggers // set-up how to log the behavior

	Services map[string]service.Service // set-up the honeypot service to catch the behavior
}

// Create a new configuration from the given path, parse and return it
func New(path string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("failed to read the configuration file")
		return nil, err
	}

	decoder := func(c *mapstructure.DecoderConfig) {
		// executes all input hook functions until one of them returns no error.
		c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			StringToURLHookFunc(),
			service.StringToServiceHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToIPNetHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		)
	}

	return config, viper.Unmarshal(&config, decoder)
}

// The customized unmarshaler for the configuration from path to Config
func (c *Config) UnmarshalText(text []byte) error {
	path := string(text)
	switch config, err := New(path); err {
	case nil:
		*c = *config
		return nil
	default:
		return err
	}
}

// Show the configuration in YAML format
func (c Config) ToYAML() string {
	switch data, err := yaml.Marshal(c); err {
	case nil:
		return string(data)
	default:
		log.Error().Err(err).Msg("failed to marshal the configuration to YAML")
		return ""
	}
}
