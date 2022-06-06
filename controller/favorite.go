package controller

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/service"
	"TikTok/vo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	var user repository.User
	redis.Get(config.UserKey+token, &user)
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if actionType == "1" {
		_, err := service.Favorite(user.ID, videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "点赞失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{StatusCode: 0, StatusMsg: "点赞成功"})
	} else if actionType == "2" {
		_, err := service.UnFavorite(user.ID, videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "取消点赞失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{StatusCode: 0, StatusMsg: "取消点赞成功"})
	}
	return
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	user, err := service.SearchUser(userId, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Response: vo.Response{StatusCode: -1, StatusMsg: "获取用户信息失败"},
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: vo.Response{StatusCode: 1, StatusMsg: "用户信息有误"},
		})
		return
	}

	videos, err := service.GetFavoriteList(userId, token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: vo.Response{StatusCode: 2, StatusMsg: "获取点赞列表失败"},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: vo.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
	return
}
