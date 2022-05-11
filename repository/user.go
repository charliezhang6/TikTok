package repository

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

type User struct {
	ID       int64  `gorm:"column:user_id"`
	Name     string `gorm:"column:user_name"`
	Password string `gorm:"column:user_password"`
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) AddUser(user User) error {
	err := db.Create(&user).Error
	if err != nil {
		log.Fatal("用户添加失败")
		return err
	}
	return nil
}

func (*UserDao) SelectByName(name string) (*User, error) {
	var user User
	err := db.Where("user_name = ?", name).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Fatal("查找用户出错" + err.Error())
	}
	return &user, err
}
