package redis

import (
	"TikTok/config"
	"fmt"
	"github.com/go-redis/redis"
)

var Client *redis.Client

func InitRedis() error {
	fmt.Println(config.Config)
	addr := config.Config.RedisAddr + ":" + config.Config.RedisPort
	fmt.Println(addr)
	Client = redis.NewClient(&redis.Options{
		Addr:     addr, // 指定
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB, // redis一共16个库，指定其中一个库即可
	})
	_, err := Client.Ping().Result()
	return err
}
