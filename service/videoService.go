package service

import (
	"TikTok/repository"
	"TikTok/vo"
	"log"
)

func Getvideos(userId int64) ([]vo.Video, int) {
	var videos []vo.Video
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

func Addvideos(video repository.Video) {
	err := repository.NewVideoDaoInstance().Addvideo(video)
	if err != nil {
		log.Println(err)
	}
}
