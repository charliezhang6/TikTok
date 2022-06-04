package service

import (
	"TikTok/repository"
	"log"
)

func GetVideos(userId int64) ([]repository.Video, int) {
	var videos []repository.Video
	videos, err := repository.NewVideoDaoInstance().SelectByUserId(userId)
	if err != nil {
		log.Println(err)
		return nil, 1
	}
	if err != nil {
		log.Println(err)
		return nil, -1
	}
	return videos, 0
}

// func returnLastvideo() ([]repository.Video,error){
// 	var videos []repository.Video
// 	err := repository.db.Limit(30).Order("date_time desc").Find(&videos).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		log.Println("查找视频出错" + err.Error())
// 	}

// 	return videos,err

// }
