package config

import (
	"errors"
	"os"
)

var AppConfig App

type App struct {
	AppKey          string
	AppName         string
	Environment     string
	AsyncmonService string
}

func loadAppConfig() error {
	appName, exist := os.LookupEnv("APP_NAME")

	if !exist {
		return errors.New("APP_NAME not set in .env")
	}

	appKey, exist := os.LookupEnv("APP_KEY")

	if !exist {
		return errors.New("APP_KEY not set in .env")
	}

	environment, exist := os.LookupEnv("ENVIRONMENT")

	if !exist {
		return errors.New("ENVIRONMENT not set in .env")
	}

	asyncmonService, exist := os.LookupEnv("ASYNCMON_SERVICE")

	if !exist {
		return errors.New("ASYNCMON_SERVICE not set in .env")
	}

	AppConfig = App{
		AppKey:          appKey,
		AppName:         appName,
		Environment:     environment,
		AsyncmonService: asyncmonService,
	}

	return nil
}
