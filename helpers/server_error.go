package helpers

import (
	"errors"
	"github.com/ichtrojan/cook/config"
)

func ServerError(err error) error {
	if config.AppConfig.Environment == "production" {
		return errors.New("something went wrong")
	}

	return err
}
