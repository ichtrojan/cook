package config

import (
	"errors"
	"os"
)

var AwsConfig AWS

type AWS struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func loadAWSConfig() error {
	awsAccessKeyID, exist := os.LookupEnv("AWS_ACCESS_KEY_ID")

	if !exist {
		return errors.New("AWS_ACCESS_KEY_ID not set in .env")
	}

	awsSecretAccessKey, exist := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	if !exist {
		return errors.New("AWS_SECRET_ACCESS_KEY not set in .env")
	}

	awsDefaultRegion, exist := os.LookupEnv("AWS_DEFAULT_REGION")

	if !exist {
		return errors.New("AWS_DEFAULT_REGION not set in .env")
	}

	AwsConfig = AWS{
		AccessKeyID:     awsAccessKeyID,
		SecretAccessKey: awsSecretAccessKey,
		Region:          awsDefaultRegion,
	}

	return nil
}
