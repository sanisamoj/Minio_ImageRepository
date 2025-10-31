package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client
var RedisOptions = &redis.Options{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisOptions = &redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	}
	RedisClient = redis.NewClient(RedisOptions)

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis connected")
}

func Set(key string, value string, expiration time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("client Redis not initialized")
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("client Redis not initialized")
	}

	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' not found", key)
	} else if err != nil {
		return "", err
	}

	return val, nil
}
