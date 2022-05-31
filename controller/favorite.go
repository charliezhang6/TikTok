package controller

import (
	"TikTok/service"
	"TikTok/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
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
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if actionType == "1" {
		_, err := service.Favorite(userId, videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "点赞失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{StatusCode: 0, StatusMsg: "点赞成功"})
	} else if actionType == "2" {
		_, err := service.UnFavorite(userId, videoId)
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
