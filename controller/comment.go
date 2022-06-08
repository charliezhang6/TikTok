package controller

import (
	"TikTok/service"
	"TikTok/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentResponse struct {
	vo.Response
	comment vo.Comment `json:"comment"`
}

type CommentListResponse struct {
	vo.Response
	CommentList []vo.Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	loginUser, err := service.GetUserByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "查询用户出错"})
		return
	}
	actionType := c.Query("action_type")
	if actionType == "1" {
		videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
		content := c.Query("comment_text")
		comment, err := service.AddComment(loginUser.ID, videoId, content)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "添加评论失败"})
			return
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: vo.Response{
				StatusCode: 0,
				StatusMsg:  "添加评论成功",
			},
			comment: *comment,
		})
	} else {
		commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		err = service.DeleteComment(commentId, loginUser.ID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "删除评论失败"})
			return
		}
		c.JSON(http.StatusOK, vo.Response{
			StatusCode: 0,
			StatusMsg:  "删除评论成功",
		})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	loginUser, err := service.GetUserByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.Response{StatusCode: -1, StatusMsg: "查询用户出错"})
		return
	}
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	comments, err := service.GetCommentList(videoId, loginUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommentListResponse{
			Response: vo.Response{
				StatusCode: -1,
				StatusMsg:  "获得评论列表失败",
			},
			CommentList: nil,
		})
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: vo.Response{
				StatusCode: 0,
				StatusMsg:  "获得评论列表成功",
			},
			CommentList: *comments,
		})
	}
}
