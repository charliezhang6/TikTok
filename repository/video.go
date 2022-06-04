package repository

import (
	"TikTok/config"
	"TikTok/vo"
	"log"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID            int64   `json:"id" gorm:"column:video_id"`                      // 视频唯一标识
	UserId        int64   `json:"-" gorm:"column:user_id"`                        // 作者唯一标识,作为联表外键，忽略json输出
	Author        vo.User `json:"author" gorm:"foreignKey:Id;references:UserId;"` // 视频作者信息
	PlayURL       string  `json:"play_url" gorm:"column:video_path"`              // 视频播放地址
	CoverURL      string  `json:"cover_url" gorm:"column:cover_path"`             // 视频封面地址
	FavoriteCount int64   `json:"favorite_count" gorm:"column:favorite_count"`    // 视频的点赞总数
	CommentCount  int64   `json:"comment_count" gorm:"column:comment_count"`      // 视频的评论总数
	//结构体变量isfavorite是否点赞等到点赞列表完成后再写
	IsFavorite bool      `json:"is_favorite" gorm:"column:is_favorite"` // true-已点赞，false-未点赞
	Title      string    `json:"title" gorm:"column:title"`             // 视频标题
	DateTime   time.Time `json:"time" gorm:"column:date_time"`          // 视频上传时间
}

type VideoDao struct {
}

var videoDao *VideoDao

// var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {

	return videoDao
}

func (*VideoDao) AddVideo(video Video) error {
	err := db.Create(&video).Error
	if err != nil {
		log.Println("视频添加失败" + err.Error())
		return err
	}
	return nil
}

func (*VideoDao) SelectByUserId(userId int64) ([]Video, error) {
	var videos []Video
	err := db.Table("videos").Joins("inner join users on videos.user_id = users.user_id where videos.user_id = ?", userId).
		Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		log.Println("查找视频出错" + err.Error())
	}
	return videos, err
}

func (*VideoDao) SelectById(videoId int64) (Video, error) {
	var video Video
	db.First(&video, "video_id = ?", videoId)
	return video, nil
}

func (*VideoDao) AddFavoriteCount(videoId int64) {
	var video Video
	db.First(&video, "video_id = ?", videoId)
	db.Model(&video).Update("favorite_count", video.FavoriteCount+1)
}

func (*VideoDao) DecrFavoriteCount(videoId int64) {
	var video Video
	db.First(&video, "video_id = ?", videoId)
	db.Model(&video).Update("favorite_count", video.FavoriteCount-1)
}

func (*VideoDao) SelectVideoById(videoID int64) (*Video, error) {
	var video Video
	err := db.Table("videos").Where("video_id = ? ", videoID).Find(&video).Error
	if err != nil {
		log.Println("查找视频出错 " + err.Error())
	}
	return &video, nil
}

func (*VideoDao) IncrCommentCount(videoId int64) {
	var video Video
	db.First(&video, "video_id = ?", videoId)
	db.Model(&video).Update("comment_count", video.CommentCount+1)
}

func (*VideoDao) DecrCommentCount(videoId int64) {
	var video Video
	db.First(&video, "video_id = ?", videoId)
	db.Model(&video).Update("comment_count", video.CommentCount-1)
}

func Feed(last_time time.Time) (video_list *[]Video, err error) {
	dbErr := db.Where("create_date <= ?", last_time.Format("2006-01-02 15:04:05")).Order("create_date desc").Limit(config.FEEDNUM).Find(&video_list)
	if dbErr.Error != nil {
		return nil, dbErr.Error
	}
	return video_list, nil
}
