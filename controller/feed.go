package controller

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/service"
	"TikTok/vo"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	vo.Response
	VideoList []repository.Video `json:"video_list,omitempty"`
	NextTime  int64              `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var latest_time time.Time
	var user_id int64
	token := c.Query("token")
	//用户不存在
	//if _, exist := usersLoginInfo[token]; !exist {
	//	c.JSON(http.StatusOK, vo.Response{StatusCode: 1,
	//		StatusMsg: "User doesn't exist"})
	//	return
	//}

	var videoUser repository.User
	err := redis.Get(config.UserKey+token, &videoUser)
	if err != nil {
		log.Println("查询redis出错" + err.Error())
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	user_id = videoUser.ID

	//不存在的时候如何创造匿名对象？

	last_time_query := c.Query("latest_time")
	if last_time_query == "" || last_time_query == "0" {
		latest_time = time.Now()
	} else {
		last_time_unix, err := strconv.ParseInt(last_time_query, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, vo.Response{
				StatusCode: 500,
				StatusMsg:  "last_time时间戳出错",
			})
			return
		}
		latestTime := time.Unix(last_time_unix/1000, 0)
		latest_time = latestTime
	}

	videoService := service.NewVideoSerVice()
	// videoService :=repository.NewVideoDaoInstance()
	resp, err := videoService.Feed(user_id, latest_time)

	if err != nil {
		c.JSON(http.StatusOK, vo.Response{
			StatusCode: 500,
			StatusMsg:  "刷新视频失败",
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}
