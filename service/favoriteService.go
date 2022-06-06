package service

import (
	"TikTok/config"
	"TikTok/redis"
	"TikTok/repository"
	"TikTok/vo"
	"log"
	"strconv"
	"time"

	redis2 "github.com/go-redis/redis"
)

func Favorite(userId int64, videoId int64) (int64, error) {
	timeStamp := time.Now().Unix()
	favorite := redis2.Z{Score: float64(timeStamp), Member: videoId}
	result, err := redis.Client.ZAdd(config.FavoriteKey+strconv.FormatInt(userId, 10), favorite).Result()
	if err != nil {
		log.Println("点赞添加失败" + err.Error())
		return 0, err
	}
	repository.NewVideoDaoInstance().AddFavoriteCount(videoId)
	return result, nil
}

func UnFavorite(userId int64, videoId int64) (int64, error) {
	result, err := redis.Client.ZRem(config.FavoriteKey+strconv.FormatInt(userId, 10), videoId).Result()
	if err != nil {
		log.Println("取消点赞失败" + err.Error())
		return 0, err
	}
	repository.NewVideoDaoInstance().DecrFavoriteCount(videoId)
	return result, nil
}

func IsFavorite(userId int64, videoId int64) (bool, error) {
	stringUserId := strconv.FormatInt(userId, 10)
	stringVideoId := strconv.FormatInt(videoId, 10)
	result, err := redis.Client.ZLexCount(config.FavoriteKey+stringUserId, "["+stringVideoId, "["+stringVideoId).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func GetFavoriteList(userId int64, token string) ([]vo.Videoinfo, error) {
	favorites, err := redis.Client.ZRange(config.FavoriteKey+strconv.FormatInt(userId, 10), 0, -1).Result()
	if err != nil {
		log.Println("获取点赞列表失败" + err.Error())
		return nil, err
	}
	videoList := make([]vo.Videoinfo, len(favorites))
	for i, favorite := range favorites {
		var videoId int64
		videoId, err = strconv.ParseInt(favorite, 10, 64)
		if err != nil {
			log.Println("类型转换失败" + err.Error())
			return nil, err
		}
		video, err := repository.NewVideoDaoInstance().SelectById(videoId)
		// video, err := repository.NewVideoDaoInstance().SelectById(1)
		if err != nil {
			return nil, err
		}
		author, err := SearchUser(video.UserId, token)
		if err != nil {
			return nil, err
		}
		video.Author = *author
		var loginUser repository.User
		err = redis.Get(config.UserKey+token, &loginUser)
		if err != nil {
			return nil, err
		}
		video.IsFavorite, _ = IsFavorite(loginUser.ID, videoId)
		//转化数据结构
		var videoinfo = &vo.Videoinfo{}
		videoinfo.Id = video.ID
		videoinfo.Author.Id = video.UserId
		videoinfo.Author = video.Author
		videoinfo.PlayUrl = video.PlayURL
		videoinfo.CoverUrl = video.CoverURL
		videoinfo.FavoriteCount = video.FavoriteCount
		videoinfo.CommentCount = video.CommentCount
		videoinfo.Title = video.Title
		videoinfo.IsFavorite, _ = IsFavorite(loginUser.ID, videoId)
		videoinfo.Author.IsFollow = author.IsFollow
		videoList[i] = *videoinfo
	}
	return videoList, nil
}
