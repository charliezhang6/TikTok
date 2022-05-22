package repository

import (
	"gorm.io/gorm"
	"log"
	"sync"
)

type User struct {
	ID          int64  `gorm:"column:user_id"`
	Name        string `gorm:"column:user_name"`
	Password    string `gorm:"column:user_password"`
	FollowCount int64  `gorm:"column:follow_count"`
	FansCount   int64  `gorm:"column:fans_count"`
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
		log.Println("用户添加失败" + err.Error())
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
		log.Println("查找用户出错" + err.Error())
	}
	return &user, err
}

func (*UserDao) AddFollowCountById(id int64) {
	var user User
	db.First(&user, "user_id=?", id)
	db.Model(&user).Update("follow_count", user.FollowCount+1)
}

func (*UserDao) AddFansCountById(id int64) {
	var user User
	db.First(&user, "user_id=?", id)
	db.Model(&user).Update("fans_count", user.FansCount+1)
}

func (*UserDao) DecrFollowCountById(id int64) {
	var user User
	db.First(&user, "user_id=?", id)
	db.Model(&user).Update("follow_count", user.FollowCount-1)
}
func (*UserDao) DecrFansCountById(id int64) {
	var user User
	db.First(&user, "user_id=?", id)
	db.Model(&user).Update("fans_count", user.FansCount-1)
}
