package redis

import (
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"user/pkg/config"
)

var redisClient *redis.Client

func InitRedis() {
	cnf := config.GetConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cnf.Redis.Host, cnf.Redis.Port),
		Password: cnf.Redis.Pwd,
	})
}
func Get() *redis.Client {
	return redisClient
}
