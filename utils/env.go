package utils

import (
	"log"
	"os"

	"github.com/klever-io/getchain-sdk/utils"
	"go.uber.org/zap"
)

var defaults = map[string]string{
	"KAFKA_BOOTSTRAP_SERVERS": "localhost:9092",
	"REDIS_CLIENT_URL":        "localhost:6379",
}

func RequireEnv(env string) string {
	str := os.Getenv(env)
	if str == "" {
		str = defaults[env]
		if str == "" {
			zap.L().Error("missing env:", zap.Error(utils.ErrMissingEnv), zap.String(env, "env"))
			log.Fatal(utils.ErrMissingEnv)
		}
	}

	return str
}
