package repository

import (
	"TikTok/vo"
	"log"

	"gorm.io/gorm"
)

//Video定义在vo/comment.go中
// type Video struct {
// 	Id            int64  `json:"id,omitempty"`
// 	Author        User   `json:"author"`
// 	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
// 	CoverUrl      string `json:"cover_url,omitempty"`
// 	FavoriteCount int64  `json:"favorite_count,omitempty"`
// 	CommentCount  int64  `json:"comment_count,omitempty"`
// 	IsFavorite    bool   `json:"is_favorite,omitempty"`
// }

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

func (*VideoDao) SelectById(userId int64) ([]vo.Video, error) {
	var videos []vo.Video
	err := db.Where("user_id = ?", userId).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("查找视频出错" + err.Error())
	}
	return videos, err
}
