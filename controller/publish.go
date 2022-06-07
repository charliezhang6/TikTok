package controller

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/service"
	"TikTok/util"
	"TikTok/vo"
	"fmt"
	"log"
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
	VideoList []vo.Videoinfo `json:"video_list"`
}

// GetCover get cover image from video file
func GetCover(filename string, filepath string) (string, string) {
	filerealname := strings.Split(filename, ".")
	filepathname := filepath + "/video/" + filename
	coverpathname := filepath + "/image/" + filerealname[0] + ".jpeg"
	println(filepathname)
	println(coverpathname)
	cmd := exec.Command("ffmpeg", "-i", filepathname, "-ss", "1", "-f", "image2", "-frames:v", "1", coverpathname)

	err := cmd.Run()
	if err != nil {
		println(err)
	}

	println("after cmd")
	return coverpathname, filepathname
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// 获取token, 不能用query
	token := c.PostForm("token")

	var videoUser repository.User
	err := redis.Get(config.UserKey+token, &videoUser)
	if err != nil {
		log.Println("查询redis出错" + err.Error())
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

	// 视频存储路径
	storePath := "/root/VideoImage"
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", videoUser.ID, filename)
	saveFile := filepath.Join(storePath+"/video/", finalName)
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

	// 获取图片保存并返回视频和图片路径
	coverpath, videopath := GetCover(finalName, storePath)
	userId := videoUser.ID
	currentTime := time.Now()
	id := util.GenSnowflake()

	// 存数据库
	video := repository.Video{
		ID:            id,
		UserId:        userId,
		DateTime:      currentTime,
		PlayURL:       videopath,
		CoverURL:      coverpath,
		Title:         title,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	service.AddVideos(video)
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//查询token和要访问的用户
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//以下为测试
	// videos, _ := service.GetVideos(1, token)
	// c.JSON(http.StatusOK, VideoListResponse{
	// 	Response:  vo.Response{StatusCode: 0},
	// 	VideoList: videos,
	// 	// VideoList: DemoTestVideos, //逗号不可省略
	// })
	// return

	//查询要访问的用户对象
	user, err := service.SearchUser(userId, token)
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
		videos, code := service.GetVideos(userId, token)
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
				// VideoList: DemoTestVideos, //逗号不可省略
			})
		}
		return
	}
}
