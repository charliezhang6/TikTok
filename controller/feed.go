package controller

import (
	"TikTok/service"
	"TikTok/vo"

	// "TikTok/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	vo.Response
	VideoList []vo.Video `json:"video_list,omitempty"`
	NextTime  int64      `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var latest_time time.Time
	var user_id int64
	token := c.Query("token")
	//用户不存在
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, vo.Response{StatusCode: 1,
			StatusMsg: "User doesn't exist"})
		return
	}

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

	videoService, k := service.GetVideos(user_id)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  vo.Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}

// // Feed same demo video list for every request
// func Feed(c *gin.Context) {
// 	//查询token和要访问的用户
// 	token := c.Query("token")
// 	if _, exist := usersLoginInfo[token]; !exist {
// 		c.JSON(http.StatusOK, vo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
// 		return
// 	}

// 	viodes,err :=service.returnLastvideo()

// 	//查询到访问的用户

// 	if err !=nil{
// 		c.JSON(http.StatusOK, FeedResponse{
// 			Response:  vo.Response{StatusCode: -1,StatusMsg: "获取用户信息失败"},
// 			VideoList: DemoVideos,
// 			NextTime:  time.Now().Unix(),
// 		})
// 	}
// 	if user==nil{
// 		c.JSON(http.StatusOK, FeedResponse{
// 			Response:  vo.Response{StatusCode: 1,StatusMsg: "用户信息有误"},
// 			VideoList: DemoVideos,
// 			NextTime:  time.Now().Unix(),
// 		})
// 	}

// 	if user !=nil{
// 		//根据信息
// 	}

// }

// func returnLastvideo() {
// 	panic("unimplemented")
// }

// Feed same demo video list for every request
// func Feed(c *gin.Context) {
// 	c.JSON(http.StatusOK, FeedResponse{
// 		Response:  vo.Response{StatusCode: 0},
// 		VideoList: DemoVideos,
// 		NextTime:  time.Now().Unix(),
// 	})
// }
