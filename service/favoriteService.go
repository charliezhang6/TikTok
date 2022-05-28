package service

import (
	"TikTok/config"
	"TikTok/redis"
	redis2 "github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

func Favorite(userId int64, videoId int64) (int64, error) {
	timeStamp := time.Now().Unix()
	favorite := redis2.Z{Score: float64(timeStamp), Member: videoId}
	result, err := redis.Client.ZAdd(config.FavoriteKey+strconv.FormatInt(userId, 10), favorite).Result()
	if err != nil {
		log.Println("点赞添加失败" + err.Error())
		return 0, err
	}
	//todo 视频表添加点赞数量
	return result, nil
}

func UnFavorite(userId int64, videoId int64) (int64, error) {
	result, err := redis.Client.ZRem(config.FollowKey+strconv.FormatInt(userId, 10), videoId).Result()
	if err != nil {
		log.Println("取消点赞失败" + err.Error())
		return 0, err
	}
	//todo 视频表减少点赞数量
	return result, nil
}
