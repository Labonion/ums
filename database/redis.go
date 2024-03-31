package database

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       int(db),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	RedisClient = client
}

func SetData(client *redis.Client, key string, value string, expiry time.Duration) error {
	err := client.Set(context.Background(), key, value, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetData(client *redis.Client, key string) (string, error) {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func DeleteKey(client *redis.Client, key string) (bool, error) {
	_, err := client.Del(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
