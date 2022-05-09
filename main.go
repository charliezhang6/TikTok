package main

import (
	"TikTok/redis"
	"TikTok/repository"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()

	initRouter(r)
	err := repository.Init()
	if err != nil {
		log.Fatal("数据库连接失败")
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	err = redis.InitRedis()
	if err != nil {
		log.Fatal("Redis连接失败")
	}
}
