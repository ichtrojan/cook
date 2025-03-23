package config

import (
	"errors"
	"os"
)

type Mail struct {
	FromAddress string
	FromName    string
}

var MailConfig Mail

func loadMailConfig() error {
	mailFrom, exist := os.LookupEnv("MAIL_FROM_ADDRESS")

	if !exist {
		return errors.New("MAIL_FROM_ADDRESS not set in .env")
	}

	MailConfig = Mail{
		FromAddress: mailFrom,
		FromName:    AppConfig.AppName,
	}

	return nil
}
