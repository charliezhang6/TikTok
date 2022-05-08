package redis

import "github.com/go-redis/redis"

var redisClient *redis.Client

func InitRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "121.40.150.207:6379", // 指定
		Password: "redis",
		DB:       0, // redis一共16个库，指定其中一个库即可
	})
	_, err := redisClient.Ping().Result()
	return err
}
