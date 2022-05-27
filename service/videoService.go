package service

import (
	"TikTok/repository"
	"log"
)

func GetVideos(userId int64) ([]repository.Video, int) {
	var videos []repository.Video
	videos, err := repository.NewVideoDaoInstance().SelectById(userId)
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
