package driver

import (
	"github.com/go-redis/redis"
)

var RedisServer *redis.Client

func NewRdids(host string) (*redis.Client, error) {
	RedisServer = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	_, err := RedisServer.Ping().Result()
	if err != nil {
		return nil, err
	}
	return RedisServer, nil
}
