package service

import (
	"TikTok/config"
	"TikTok/redis"
	redis2 "github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

func Follow(fromId int64, toId int64) (int64, error) {
	timeStamp := time.Now().Unix()
	fans := redis2.Z{Score: float64(timeStamp), Member: toId}
	result, err := redis.Client.ZAdd(config.FollowKey+strconv.FormatInt(fromId, 10), fans).Result()
	if err != nil {
		log.Println("添加关注表失败" + err.Error())
		return 0, err
	}
	return result, nil
}

//func UnFollow(fromId int64, toId int64) {
//
//}
