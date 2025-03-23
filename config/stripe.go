package config

import (
	"errors"
	"os"
)

var StripeConfig Stripe

type Stripe struct {
	SecretKey     string
	WebhookSecret string
}

func loadStripeConfig() error {
	stripeSecretKey, exist := os.LookupEnv("STRIPE_SECRET_KEY")

	if !exist {
		return errors.New("STRIPE_SECRET_KEY not set in .env")
	}

	stripeWebhookSecret, exist := os.LookupEnv("STRIPE_WEBHOOK_SECRET")

	if !exist {
		return errors.New("STRIPE_WEBHOOK_SECRET not set in .env")
	}

	StripeConfig = Stripe{
		SecretKey:     stripeSecretKey,
		WebhookSecret: stripeWebhookSecret,
	}

	return nil
}
