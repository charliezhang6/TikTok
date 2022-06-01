package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"log"
)

func GetVideos(userId int64, token string) ([]repository.Video, int) {
	var videos []repository.Video
	//拉取视频列表数据（无关注和点赞）
	videos, err := repository.NewVideoDaoInstance().SelectByUserId(userId)
	if err != nil {
		log.Println(err)
		return nil, 1
	}
	if err != nil {
		log.Println(err)
		return nil, -1
	}
	//循环处理视频列表中的关注和点赞
	var loginUser repository.User
	err = redis.Get(config.UserKey+token, &loginUser) //获取当前登录用户信息
	if err != nil {
		return nil, 1
	}
	for i, videoinfo := range videos {
		//获取视频ID和作者ID
		videoId := videoinfo.ID
		AuthorId := videoinfo.Author.Id
		//获得视频点赞信息
		videos[i].IsFavorite, _ = IsFavorite(loginUser.ID, videoId)
		//获得作者关注信息
		author, err := SearchUser(AuthorId, token)
		if err != nil {
			return nil, 1
		}
		videos[i].Author.IsFollow = author.IsFollow
	}
	//返回指定信息
	return videos, 0
}
