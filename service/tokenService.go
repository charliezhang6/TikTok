package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
)

func CreateToken(user *repository.User) string {
	token := user.Name + user.Password
	return token
}

func RefreshToken(token string, user *repository.User) error {
	err := redis.Set(config.LoginTokenKey+token, user, config.Config.ExpireTime*60)
	return err
}
