package queue

import (
	"crypto/tls"
	"github.com/hibiken/asynq"
	"github.com/ichtrojan/cook/config"
	"github.com/ichtrojan/cook/mailer"
)

var Client *asynq.Client

func Register() *asynq.ServeMux {
	redisConnection := asynq.RedisClientOpt{
		Addr:     config.RedisConfig.Addr,
		Username: config.RedisConfig.User,
		Password: config.RedisConfig.Password,
	}

	if config.RedisConfig.Scheme == "tls" {
		redisConnection.TLSConfig = &tls.Config{}
	}

	Client = asynq.NewClient(redisConnection)

	mux := asynq.NewServeMux()

	mux.HandleFunc("send:email", mailer.Send)

	return mux
}
