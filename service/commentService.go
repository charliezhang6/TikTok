package service

import (
	"TikTok/repository"
	"TikTok/util"
	"TikTok/vo"
	"errors"
	"log"
	"strconv"
	"time"
)

func AddComment(userId int64, videoId int64, content string) (*vo.Comment, error) {
	user, err := repository.NewUserDaoInstance().SelectById(userId)
	if err != nil || user.ID == 0 {
		return nil, errors.New("查找数据库失败，userId或许不存在")
	}
	video, err := repository.NewVideoDaoInstance().SelectVideoById(videoId)
	if err != nil || video.ID == 0 {
		return nil, errors.New("查找数据库失败，videoId或许不存在")
	}
	comment := &repository.Comment{
		ID:         util.GenSnowflake(),
		UserId:     userId,
		VideoId:    videoId,
		Content:    content,
		CreateDate: time.Now().Month().String() + "-" + strconv.Itoa(time.Now().Day()),
	}
	err = repository.NewCommentDaoInstance().AddComment(comment)
	if err != nil {
		return nil, err
	}
	repository.NewVideoDaoInstance().IncrCommentCount(videoId)
	voComment, err := fillVoComment(user, comment, userId)
	if err != nil {
		return nil, err
	}
	return voComment, nil
}

func DeleteComment(commentId int64, userId int64) error {
	var err error
	comment, err := repository.NewCommentDaoInstance().SelectCommentById(commentId)
	if err != nil || comment.ID == 0 {
		return errors.New("查找数据库失败，comment或许不存在")
	}
	if comment.UserId != userId {
		return errors.New("没有权限删除")
	}
	err = repository.NewCommentDaoInstance().DeleteComment(comment)
	if err != nil {
		return err
	}
	repository.NewVideoDaoInstance().DecrCommentCount(comment.VideoId)
	return nil
}

func GetCommentList(videoId int64, loginUserId int64) (*[]vo.Comment, error) {
	comments, err := repository.NewCommentDaoInstance().SelectCommentsByVideoId(videoId)
	if err != nil {
		return nil, err
	}
	voComments := make([]vo.Comment, 0)
	for _, comment := range *comments {
		voComment, err := generateVoComment(&comment, loginUserId)
		if err != nil {
			log.Println("获得comment列表失败 " + err.Error())
			return nil, err
		}
		voComments = append(voComments, *voComment)
	}
	return &voComments, nil
}

func generateVoComment(comment *repository.Comment, loginUserId int64) (*vo.Comment, error) {
	user, err := repository.NewUserDaoInstance().SelectById(comment.UserId)
	if err != nil {
		return nil, err
	}
	voComments, err := fillVoComment(user, comment, loginUserId)
	if err != nil {
		return nil, err
	}
	return voComments, nil
}

func fillVoComment(user *repository.User, comment *repository.Comment, loginUserId int64) (*vo.Comment, error) {
	isFollow, err := IsFollow(loginUserId, user.ID)
	if err != nil {
		return nil, err
	}
	voUser := vo.User{
		Id:            user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow:      isFollow,
	}

	return &vo.Comment{
		Id:         comment.ID,
		User:       voUser,
		Content:    comment.Content,
		CreateDate: comment.CreateDate,
	}, nil
}
