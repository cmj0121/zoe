package config

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Loggers []Logger

// Run the loggers to receive the message from community channel and
// send to the specified logger implemented
func (l *Loggers) Run(ch <-chan *types.Message) error {
	log.Debug().Msg("starting the logger ...")
	defer log.Debug().Msg("finished the logger ...")

	if err := l.prologue(); err != nil {
		log.Error().Err(err).Msg("failed to set up the logger")
		return err
	}

	return l.run(ch)
}

// setup the logger before running the logger
func (l *Loggers) prologue() error {
	for index, logger := range *l {
		switch logger.Scheme {
		case "console", "stdout":
			(*l)[index].WriteCloser = &ConsoleLogger{}
		case "file":
			path := fmt.Sprintf("%s%s", logger.Host, logger.Path)
			file, err := NewFileLogger(path)
			if err != nil {
				log.Error().Err(err).Str("path", path).Msg("failed to open the file")
				return err
			}

			(*l)[index].WriteCloser = file
		case "sqlite", "sqlite3":
			path := fmt.Sprintf("%s%s", logger.Host, logger.Path)
			db, err := NewSQLiteLogger(path)
			if err != nil {
				log.Error().Err(err).Str("path", path).Msg("failed to open the sqlite")
				return err
			}

			(*l)[index].WriteCloser = db
		default:
			err := fmt.Errorf("unknown logger: %s", logger.String())
			return err
		}
	}

	return nil
}

// run the logger to receive the message from the community channel
func (l *Loggers) run(ch <-chan *types.Message) error {
	for {
		select {
		case msg, ok := <-ch:
			log.Info().Interface("message", msg).Interface("message", msg).Msg("received a message")
			if !ok {
				log.Debug().Msg("no message")
				break
			}

			allErrors := []error{}
			for _, logger := range *l {
				switch logger.WriteCloser {
				case nil:
					log.Warn().Str("logger", logger.String()).Msg("no writer")
				default:
					if err := logger.WriteCloser.Write(msg); err != nil {
						log.Warn().Err(err).Str("logger", logger.String()).Msg("cannot handle message")
						allErrors = append(allErrors, err)
					}
				}
			}

			err := errors.Join(allErrors...)
			if err != nil {
				log.Error().Err(err).Msg("failed to handle the message")
				continue
			}
		}
	}

	return nil
}

type Logger struct {
	// The final URL to store the message
	*url.URL

	// The writer to store the message and close the connection
	types.WriteCloser
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
