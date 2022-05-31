package repository

import (
	"log"
	"sync"
)

type Comment struct {
	ID         int64  `gorm:"column:comment_id"`
	UserId     int64  `gorm:"column:user_id"`
	VideoId    int64  `gorm:"column:video_id"`
	Content    string `gorm:"column:content"`
	CreateDate string `gorm:"column:create_date"`
}

type CommentDao struct{}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

func (*CommentDao) SelectCommentById(commentId int64) (*Comment, error) {
	var comment Comment
	err := db.Table("comments").Where("comment_id = ?", commentId).First(&comment).Error
	if err != nil {
		log.Println("通过ID查找评论失败 " + err.Error())
	}
	return &comment, err
}

func (*CommentDao) AddComment(comment *Comment) error {
	err := db.Create(&comment).Error
	if err != nil {
		log.Println("评论添加失败 " + err.Error())
	}
	return err
}

func (*CommentDao) DeleteComment(comment *Comment) error {
	err := db.Delete(comment).Error
	if err != nil {
		log.Println("评论删除失败 " + err.Error())
	}
	return err
}

func (*CommentDao) SelectCommentsByVideoId(videoId int64) (*[]Comment, error) {
	var comments []Comment
	err := db.Table("comments").Joins("inner join videos on comments.video_id = videos.video_id"+
		" where comments.video_id = ?", videoId).Find(&comments).Error
	if err != nil {
		log.Println("查找评论出错" + err.Error())
	}
	return &comments, err
}
