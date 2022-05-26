package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/vo"
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

func GetFollowList(id int64) ([]vo.User, error) {
	followings, err := redis.Client.ZRange(config.FollowKey+strconv.FormatInt(id, 10), 0, -1).Result()
	if err != nil {
		log.Println("获取关注列表失败" + err.Error())
		return nil, err
	}
	userList := make([]vo.User, len(followings))
	for i, following := range followings {
		var toId int64
		toId, err = strconv.ParseInt(following, 10, 64)
		if err != nil {
			log.Println("类型转换失败" + err.Error())
			return nil, err
		}
		user, err := repository.NewUserDaoInstance().SelectById(toId)
		if err != nil {
			return nil, err
		}
		userList[i] = vo.User{
			Id:            toId,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow:      true,
		}
	}
	return userList, err
}

func GetFanList(id int64) ([]vo.User, error) {
	sid := strconv.FormatInt(id, 10)
	fans, err := redis.Client.ZRange(config.FansKey+sid, 0, -1).Result()
	if err != nil {
		log.Println("获取关注列表失败" + err.Error())
		return nil, err
	}
	userList := make([]vo.User, len(fans))
	for i, fan := range fans {
		var toId int64
		toId, err := strconv.ParseInt(fan, 10, 64)
		if err != nil {
			log.Println("类型转换失败" + err.Error())
			return nil, err
		}
		user, err := repository.NewUserDaoInstance().SelectById(toId)
		if err != nil {
			return nil, err
		}
		isFollow, err := IsFollow(id, user.ID)
		userList[i] = vo.User{
			Id:            toId,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow:      isFollow,
		}
	}
	return userList, nil
}

func IsFollow(fromId int64, toId int64) (bool, error) {
	stringFromId := strconv.FormatInt(fromId, 10)
	stringToId := strconv.FormatInt(toId, 10)
	result, err := redis.Client.ZLexCount(config.FollowKey+stringFromId, "["+stringToId, "["+stringToId).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
