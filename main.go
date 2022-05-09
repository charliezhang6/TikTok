package main

import (
	"TikTok/controller"
	"TikTok/redis"
	"fmt"
	"log"
)

func main() {
	//r := gin.Default()
	//
	//initRouter(r)
	//err := repository.Init()
	//if err != nil {
	//	log.Fatal("数据库连接失败")
	//}
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	err := redis.InitRedis()
	if err != nil {
		log.Fatal("Redis连接失败")
	}
	var user controller.User = controller.User{Name: "test", Id: 1111, FollowCount: 123, FollowerCount: 12345, IsFollow: true}
	redis.Set("testtoken", user, 0)

	//jsonbyte, _ := util.DefaultTranscoder.Marshal(user)
	var newuser controller.User
	redis.Get("testtoken", &newuser)
	fmt.Println(newuser)
}
