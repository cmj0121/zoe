package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Loggers []string

// The configuration of the Zoe service including
//   - how to log the behavior
//   - the web interface to show the behavior
//   - the honeypot service to catch the behavior
type Config struct {
	Loggers // set-up how to log the behavior
}

func New(path string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("failed to read the configuration file")
		return nil, err
	}

	return config, viper.Unmarshal(&config)
}

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

func (c Config) ToYAML() string {
	switch data, err := yaml.Marshal(c); err {
	case nil:
		return string(data)
	default:
		log.Error().Err(err).Msg("failed to marshal the configuration to YAML")
		return ""
	}
}
