package repository

import (
	"TikTok/config"
	"TikTok/vo"
	"log"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID            int64   `gorm:"column:video_id" json:"id"`                      // 视频唯一标识
	UserId        int64   `gorm:"column:user_id" json:"-"`                        // 作者唯一标识,作为联表外键，忽略json输出
	Author        vo.User `gorm:"foreignKey:Id;references:UserId;" json:"author"` // 视频作者信息
	PlayURL       string  `gorm:"column:video_path" json:"play_url"`              // 视频播放地址
	CoverURL      string  `gorm:"column:cover_path" json:"cover_url"`             // 视频封面地址
	FavoriteCount int64   `gorm:"column:favorite_count" json:"favorite_count"`    // 视频的点赞总数
	CommentCount  int64   `gorm:"column:comment_count" json:"comment_count"`      // 视频的评论总数
	//结构体变量isfavorite是否点赞等到点赞列表完成后再写
	IsFavorite bool      `json:"is_favorite"`                  // true-已点赞，false-未点赞
	Title      string    `gorm:"column:title" json:"title"`    // 视频标题
	DateTime   time.Time `json:"time" gorm:"column:date_time"` // 视频上传时间
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
	//注：剩下要解决的bug是没有在video.author的vo.user中加入gorm标签，所以无法把author拉过来，这个之后再改
	//author拉过来之后是初始值，不是数据库里面的值，问题在哪里呢
	//要使用preload函数提前加载users，则问题解决
	err := db.Table("videos").Preload("Author").Joins("inner join users on videos.user_id = users.user_id").Where("videos.user_id = ?", userId).
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
	dbErr := db.Where("date_time <= ?", last_time.Format("2006-01-02 15:04:05")).Order("date_time desc").Limit(config.FEEDNUM).Find(&video_list)
	if dbErr.Error != nil {
		return nil, dbErr.Error
	}
	return video_list, nil
}
