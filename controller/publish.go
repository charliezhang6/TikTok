package controller

import (
	"TikTok/service"
	"TikTok/vo"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	vo.Response
	VideoList []vo.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, vo.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, vo.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//查询token和鉴权id
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	//判断用户与token是否相同并从用户表中获取用户信息句柄
	user, err := service.CheckUser(userId, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, VideoListResponse{
			Response: vo.Response{StatusCode: -1, StatusMsg: "获取用户信息失败"},
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: vo.Response{StatusCode: 1, StatusMsg: "用户信息有误"},
		})
		return
	}
	if user != nil {
		//根据用户信息从视频表中获取视频信息
		videos, code := service.Getvideos(userId)
		//返回结果json
		if code != 0 { //尚未发布视频
			c.JSON(http.StatusOK, VideoListResponse{
				Response:  vo.Response{StatusCode: 0},
				VideoList: nil,
			})
		}
		if code == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response:  vo.Response{StatusCode: 0},
				VideoList: videos,
				// VideoList: DemoVideos, //逗号不可省略
			})
		}
		return
	}
}
