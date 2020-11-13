package common

import (
	"log"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedisClient(addr string, dbName int, poolSize int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       dbName,
		PoolSize: 10,
	})

	log.Println("Connected to Redis!")

	RDB = cli
	return cli
}
