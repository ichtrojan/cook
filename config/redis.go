package config

import (
	"errors"
	"fmt"
	"os"
)

var RedisConfig Redis

type Redis struct {
	Host     string
	User     string
	Port     string
	Password string
	Scheme   string
	Addr     string
}

func loadRedisConfig() error {
	redisHost, exist := os.LookupEnv("REDIS_HOST")

	if !exist {
		return errors.New("REDIS_HOST not set in .env")
	}

	redisPort, exist := os.LookupEnv("REDIS_PORT")

	if !exist {
		return errors.New("REDIS_PORT not set in .env")
	}

	redisUser, exist := os.LookupEnv("REDIS_USER")

	if !exist {
		return errors.New("REDIS_USER not set in .env")
	}

	redisPass, exist := os.LookupEnv("REDIS_PASS")

	if !exist {
		return errors.New("REDIS_PASS not set in .env")
	}

	redisScheme, exist := os.LookupEnv("REDIS_SCHEME")

	if !exist {
		return errors.New("REDIS_SCHEME not set in .env")
	}

	RedisConfig = Redis{
		Host:     redisHost,
		User:     redisUser,
		Port:     redisPort,
		Password: redisPass,
		Scheme:   redisScheme,
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
	}

	return nil
}
