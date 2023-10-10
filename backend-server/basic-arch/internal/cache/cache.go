package cache

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

var (
	redisClient *redis.Client
)

func NewCache() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PWD"),            // no password set
		DB:       cast.ToInt(os.Getenv("REDIS_DB")), // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis ping failed", err)
	}
	redisClient = rdb
}

func GetRdbClient() *redis.Client {
	return redisClient
}

func GetValueFromRedis(ctx context.Context, key string, defaultValue string) (string, error) {
	val, err := GetRdbClient().Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return defaultValue, nil
		}
		return "", err
	}
	return val, err
}
