package controller

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/service"
	"TikTok/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	vo.Response
	UserList []vo.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	var user repository.User
	redis.Get(config.UserKey+token, &user)
	actionType := c.Query("action_type")
	targetId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if actionType == "1" {
		_, err := service.Follow(user.ID, targetId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "关注失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{StatusCode: 0, StatusMsg: "关注成功"})
	} else if actionType == "2" {
		_, err := service.UnFollow(user.ID, targetId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "取关失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{StatusCode: 0, StatusMsg: "取关成功"})
	}
	return
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	user, err := service.CheckUser(userId, token)

	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "查询用户出错"})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "用户信息有误"})
		return
	}

	userList, err := service.GetFollowList(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "获取关注列表失败"})
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: vo.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	user, err := service.CheckUser(userId, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "查询用户出错"})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "用户信息有误"})
		return
	}

	userList, err := service.GetFanList(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "获取粉丝列表失败"})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: vo.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
