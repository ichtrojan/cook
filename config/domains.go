package config

import (
	"errors"
	"os"
)

var DomainConfig Domain

type Domain struct {
	ApiHost      string
	FrontendHost string
}

func loadDomainConfig() error {
	apiHost, exist := os.LookupEnv("API_HOST")

	if !exist {
		return errors.New("API_HOST not set in .env")
	}

	frontendHost, exist := os.LookupEnv("FRONTEND_HOST")

	if !exist {
		return errors.New("FRONTEND_HOST not set in .env")
	}

	DomainConfig = Domain{
		ApiHost:      apiHost,
		FrontendHost: frontendHost,
	}

	return nil
}
