package repository

import (
	"TikTok/vo"
	"log"

	"gorm.io/gorm"
)

// Video
type Video struct {
	ID            int64   `json:"id"`             // 视频唯一标识
	Author        vo.User `json:"author"`         // 视频作者信息
	PlayURL       string  `json:"play_url"`       // 视频播放地址
	CoverURL      string  `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64   `json:"favorite_count"` // 视频的点赞总数
	CommentCount  int64   `json:"comment_count"`  // 视频的评论总数
	IsFavorite    bool    `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title         string  `json:"title"`          // 视频标题
}

type VideoDao struct {
}

var videoDao *VideoDao

// var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {

	return videoDao
}

// func (*VideoDao) Addvideo(user User) error {

// 	return nil
// }

func (*VideoDao) SelectById(userId int64) ([]Video, error) {
	var videos []Video
	err := db.Where("user_id = ?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("查找视频出错" + err.Error())
	}
	return videos, err
}
