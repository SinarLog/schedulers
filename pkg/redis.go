package pkg

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func GetRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		ClientName: "sinarlog",
		Password:   os.Getenv("REDIS_PASSWORD"),
		Username:   os.Getenv("REDIS_USERNAME"),
		DB:         0,
	})
}
