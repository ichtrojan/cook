package config

import (
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() error {
	_ = godotenv.Load()

	if err := loadAppConfig(); err != nil {
		return err
	}

	if err := LoadMySQLConfig(); err != nil {
		return err
	}

	if err := loadAWSConfig(); err != nil {
		return err
	}

	if err := loadMailConfig(); err != nil {
		return err
	}

	if err := loadDomainConfig(); err != nil {
		return err
	}

	if err := loadRedisConfig(); err != nil {
		return err
	}

	if err := loadStripeConfig(); err != nil {
		return err
	}

	return nil
}
