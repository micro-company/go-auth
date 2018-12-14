package redis

import (
	"github.com/go-redis/redis"
	"github.com/micro-company/go-auth/utils"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	Redis *redis.Client
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToRedis() {
	// Get configuration
	REDIS_URL := utils.Getenv("REDIS_URL", "redis://localhost:6379/1")
	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)

	_, err = Redis.Ping().Result()
	if err != nil {
		log.Panic("Fail connect to Redis")
		panic(err)
	}

	log.Info("Success connect to Redis")
}
