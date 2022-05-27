package main

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config.InitConfig()
	r := gin.Default()

	initRouter(r)
	err := repository.Init()
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	err = redis.InitRedis()
	if err != nil {
		log.Fatal("Redis连接失败")
	}
	r.Run("0.0.0.0:9898")
}
