package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/vo"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type VideoService struct {
}

func NewVideoSerVice() *VideoService {
	return &VideoService{}
}

func GetVideos(userId int64, token string) ([]vo.Video, int) {
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
	var publishvideos []vo.Video //返回的发布列表
	for _, video := range videos {
		//获取视频ID和作者ID
		videoId := video.ID
		AuthorId := video.Author.Id
		var videoinfo = &vo.Video{}
		videoinfo.Id = video.ID
		videoinfo.Author.Id = video.UserId
		videoinfo.Author = video.Author
		videoinfo.PlayUrl = video.PlayURL
		videoinfo.CoverUrl = video.CoverURL
		videoinfo.FavoriteCount = video.FavoriteCount
		videoinfo.CommentCount = video.CommentCount
		videoinfo.Title = video.Title
		//点赞信息
		videoinfo.IsFavorite, _ = IsFavorite(loginUser.ID, videoId)
		// videoinfo.IsFavorite = true
		//获得作者关注信息
		author, err := SearchUser(AuthorId, token)
		if err != nil {
			return nil, 1
		}
		videoinfo.Author.IsFollow = author.IsFollow
		publishvideos = append(publishvideos, *videoinfo)
	}
	//返回指定信息
	return publishvideos, 0
}

func AddVideos(video repository.Video) {
	err := repository.NewVideoDaoInstance().AddVideo(video)
	if err != nil {
		log.Println(err)
	}
}

type FeedResponse struct {
	vo.Response
	NextTime  int64               `json:"next_time"`
	VideoList *[]repository.Video `json:"video_list"`
}

func (vs *VideoService) Feed(user_id int64, last_time time.Time) (resp *FeedResponse, err error) {
	var next_time time.Time
	var videoList = make([]repository.Video, 0, config.FEEDNUM)
	videos, err := repository.Feed(last_time)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, n := 0, len(*videos); i < n; i++ {
		var videoDao_t = (*videos)[i]
		wg.Add(1)
		go func(videoDao repository.Video) {
			defer wg.Done()
			var video = &repository.Video{}
			video.ID = videoDao.ID
			video.UserId = videoDao.UserId
			video.Author = videoDao.Author
			video.PlayURL = videoDao.PlayURL
			video.CoverURL = videoDao.CoverURL
			video.FavoriteCount = videoDao.FavoriteCount
			video.CommentCount = videoDao.CommentCount
			videoList = append(videoList, *video)
		}(videoDao_t)
	}

	if len(*videos) < config.FEEDNUM {
		next_time = time.Now()
	} else {
		next_time = (*videos)[len(*videos)-1].DateTime
	}
	wg.Wait()
	sort.Slice(videoList, func(i, j int) bool {
		return videoList[i].DateTime.After(videoList[j].DateTime)
	})
	return &FeedResponse{
		Response: vo.Response{
			StatusCode: 200,
			StatusMsg:  fmt.Sprintf("刷新%d条视频", len(*videos)),
		},
		NextTime:  next_time.Unix(),
		VideoList: &videoList,
	}, nil
}
