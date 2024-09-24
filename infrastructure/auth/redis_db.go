package auth

import (
	"github.com/go-redis/redis/v7"
	"os"
)

type RedisService struct {
	Auth   AuthInterface
	Client *redis.Client
}

func NewRedisDB() (*RedisService, error) {
	//redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
	return &RedisService{
		Auth:   NewAuth(redisClient),
		Client: redisClient,
	}, nil
}
