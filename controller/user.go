package controller

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if exist, _ := repository.NewUserDaoInstance().SelectByName(username); exist != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		userId, token, err := service.Register(username, password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserLoginResponse{
				Response: Response{StatusCode: -1, StatusMsg: "用户注册失败"},
			})
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userId,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	code, loginUser := service.Login(username, password)
	if code == -1 {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: -1, StatusMsg: "Fail to login"},
		})
	}

	if code == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Username or password is wrong"},
		})
	}

	if code == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   loginUser.User.ID,
			Token:    loginUser.Token,
		})
	}

}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	var user User
	err := redis.Get(config.UserKey+token, &user)
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Response: Response{StatusCode: -1, StatusMsg: "获取用户信息失败"},
		})
		return
	}
	if user.Id != userId {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: -1, StatusMsg: "用户信息有误"},
		})
	}
	if user.Id != 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
