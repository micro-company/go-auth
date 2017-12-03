package db

import (
	"github.com/micro-company/go-auth/utils"
	"github.com/go-redis/redis"
)

var (
	Redis *redis.Client
)

func ConnectToRedis() {
	// Get configuration
	REDIS_URL := utils.Getenv("REDIS_URL", "redis://localhost:6379")
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
