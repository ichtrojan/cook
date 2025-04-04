package queue

import (
	"crypto/tls"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/ichtrojan/cook/config"
	"github.com/ichtrojan/cook/mailer"
)

var Client *asynq.Client

func Register() *asynq.ServeMux {
	redisAddr := fmt.Sprintf("%s:%s", config.RedisConfig.Host, config.RedisConfig.Port)

	redisConnection := asynq.RedisClientOpt{
		Addr:     redisAddr,
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
