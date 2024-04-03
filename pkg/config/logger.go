package config

import (
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Loggers []Logger

type Logger struct {
	*url.URL
}

func (l Logger) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// Convert the string to a URL
func StringToURLHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			log.Debug().Str("from", f.String()).Msg("not a string")
			return data, nil
		}

		switch t {
		case reflect.TypeOf(url.URL{}):
			return url.Parse(data.(string))
		case reflect.TypeOf(Logger{}):
			url, err := url.Parse(data.(string))
			if err != nil {
				return data, err
			}
			return Logger{URL: url}, nil
		default:
			log.Info().Str("to", t.String()).Msg("not a URL")
			return data, nil
		}
	}
}
