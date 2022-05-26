package controller

import (
	"TikTok/repository"
	"TikTok/service"
	"TikTok/util"
	"TikTok/vo"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	vo.Response
	VideoList []vo.Video `json:"video_list"`
}

// GetCover get cover image from video file
func GetCover(filename string, filepath string) (string, string) {
	filerealname := strings.Split(filename, ".")
	filepathname := filepath + "/video/" + filename
	coverpathname := filepath + "/image/" + filerealname[0] + ".jpeg"
	cmd := exec.Command("ffmpeg", "-i", filepathname, "-ss", "1", "-f", "image2", coverpathname)

	cmd.Run()
	return coverpathname, filepathname
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	title := c.PostForm("title")
	if err != nil {
		c.JSON(http.StatusOK, vo.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	storePath := "D:/gotmp/"
	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join(storePath, finalName)
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
	coverpath, videopath := GetCover(finalName, storePath)
	userId := usersLoginInfo[token].Id
	currentTime := time.Now()
	id := util.GenSonyflake()
	video := repository.Video{
		VideoId:   id,
		UserId:    userId,
		DateTime:  currentTime,
		VideoPath: videopath,
		CoverPath: coverpath,
		Title:     title,
	}
	service.Addvideos(video)
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
