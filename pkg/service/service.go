package service

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/service/ssh"
	"github.com/cmj0121/zoe/pkg/service/types"
)

// The virtual service interface to be implemented the honeypot service
type Service interface {
	// execute the service and return the error if any
	Run(ch chan<- *types.Message) error
}

func StringToServiceHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			log.Debug().Str("from", f.String()).Msg("not a map")
			return data, nil
		}

		switch t {
		case reflect.TypeOf(map[string]Service{}):
			data := data.(map[string]interface{})
			return mapToService(data)
		default:
			log.Debug().Str("to", t.String()).Msg("not a map[string]Service")
			return data, nil
		}
	}
}

func mapToService(data map[string]interface{}) (interface{}, error) {
	svc := map[string]Service{}

	for key, value := range data {
		switch key {
		case ssh.SVC_NAME:
			sshService := ssh.SSH{}

			switch err := mapstructure.Decode(value, &sshService); err {
			case nil:
				svc[key] = &sshService
			default:
				log.Error().Err(err).Str("service", key).Msg("failed to decode the service")
				return data, err
			}
		default:
			log.Info().Str("service", key).Msg("unknown service name")
			return data, fmt.Errorf("unknown service name: %s", key)
		}
	}

	return svc, nil
}
