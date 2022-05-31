package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"crypto/rand"
	"encoding/base32"
	"log"
)

func CreateToken() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("生成token出错" + err.Error())
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return token
}

func RefreshToken(token string, user *repository.User) error {
	err := redis.Set(config.UserKey+token, user, config.Config.ExpireTime*60)
	return err
}

func GetUserByToken(token string) (*repository.User, error) {
	var user repository.User
	var err error
	err = redis.Get(config.UserKey+token, &user)
	if err != nil {
		log.Println("查询redis出错" + err.Error())
		return nil, err
	}
	err = RefreshToken(token, &user)
	if err != nil {
		log.Println("刷新token失败" + err.Error())
		return nil, err
	}
	return &user, nil
}
