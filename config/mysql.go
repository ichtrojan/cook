package config

import (
	"errors"
	"os"
)

var MySQLConfig Mysql

type Mysql struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func LoadMySQLConfig() error {
	dbHost, exist := os.LookupEnv("DB_HOST")

	if !exist {
		return errors.New("DB_HOST not set in .env")
	}

	dbPort, exist := os.LookupEnv("DB_PORT")

	if !exist {
		return errors.New("DB_PORT not set in .env")
	}

	dbUser, exist := os.LookupEnv("DB_USER")

	if !exist {
		return errors.New("DB_USER not set in .env")
	}

	dbPass, exist := os.LookupEnv("DB_PASS")

	if !exist {
		return errors.New("DB_PASS not set in .env")
	}

	dbName, exist := os.LookupEnv("DB_NAME")

	if !exist {
		return errors.New("DB_NAME not set in .env")
	}

	MySQLConfig = Mysql{
		Host: dbHost,
		Port: dbPort,
		User: dbUser,
		Pass: dbPass,
		Name: dbName,
	}

	return nil
}
