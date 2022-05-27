package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/util"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type LoginUser struct {
	User  repository.User
	Token string
}

func Register(name string, password string) (int64, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	userId := util.GenSnowflake()
	user := repository.User{ID: userId, Name: name, Password: string(hashedPassword)}
	err = repository.NewUserDaoInstance().AddUser(user)
	if err != nil {
		log.Println("添加用户出错" + err.Error())
		return 0, "", err
	}
	token := CreateToken()
	err = RefreshToken(token, &user)
	if err != nil {
		log.Println("添加token出错" + err.Error())
		return 0, "", err
	}
	return userId, token, nil
}

func Login(name string, password string) (int, *LoginUser) {
	code, user := authenticateUser(name, password)
	if code != 0 {
		return code, nil
	}
	token := CreateToken()
	if err := RefreshToken(token, user); err != nil {
		return -1, nil
	}
	return 0, &LoginUser{
		User:  *user,
		Token: token,
	}
}

// authenticateUser Returns status code
func authenticateUser(name string, password string) (int, *repository.User) {
	user, err := repository.NewUserDaoInstance().SelectByName(name)
	if err != nil {
		log.Println(err)
		return 1, nil
	}
	if err != nil {
		log.Println(err)
		return -1, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return 1, nil
	}
	return 0, user
}

func CheckUser(userId int64, token string) (*repository.User, error) {
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
	if userId == user.ID {
		return &user, nil
	}
	return nil, nil
}
