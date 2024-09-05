package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

func NewRedis() *redis.Client {

	dbInt, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbInt,
	})

	return rdb
}
