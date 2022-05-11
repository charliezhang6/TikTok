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
		log.Fatal("生成token出错" + err.Error())
	}
	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	return token
}

func RefreshToken(token string, user *repository.User) error {
	err := redis.Set(config.LoginTokenKey+token, user, config.Config.ExpireTime*60)
	return err
}
