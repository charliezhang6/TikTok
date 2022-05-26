package repository

import (
	"TikTok/vo"
	"log"
	"time"

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

type Video struct {
	VideoId   int64     `gorm:"column:video_id"`
	UserId    int64     `gorm:"column:user_id"`
	DateTime  time.Time `gorm:"column:date_time"`
	VideoPath string    `gorm:"column:video_path"`
	CoverPath string    `gorm:"column:cover_path"`
	Title     string    `gorm:"column:title"`
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

func (*VideoDao) Addvideo(video Video) error {
	err := db.Create(&video).Error
	if err != nil {
		log.Println("视频添加失败" + err.Error())
		return err
	}
	return nil
}
