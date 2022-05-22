package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
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
	repository.NewUserDaoInstance().AddFollowCountById(fromId)
	follow := redis2.Z{Score: float64(timeStamp), Member: fromId}
	result, err = redis.Client.ZAdd(config.FansKey+strconv.FormatInt(toId, 10), follow).Result()
	if err != nil {
		log.Println("添加粉丝表失败" + err.Error())
		return 0, err
	}
	repository.NewUserDaoInstance().AddFansCountById(toId)
	return result, nil
}

func UnFollow(fromId int64, toId int64) (int64, error) {
	result, err := redis.Client.ZRem(config.FollowKey+strconv.FormatInt(fromId, 10), toId).Result()
	if err != nil {
		log.Println("删除关注表失败" + err.Error())
		return 0, err
	}
	repository.NewUserDaoInstance().DecrFollowCountById(fromId)
	result, err = redis.Client.ZRem(config.FansKey+strconv.FormatInt(toId, 10), fromId).Result()
	if err != nil {
		log.Println("删除粉丝表失败" + err.Error())
		return 0, err
	}
	repository.NewUserDaoInstance().DecrFansCountById(toId)
	return result, nil
}
