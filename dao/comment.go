package dao

import (
	"gorm.io/gorm"
	"sync"
)

type Comment struct {
	Id         int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 主键id
	VideoID    int    `gorm:"column:video_id" json:"-"`
	UserID     int    `gorm:"column:user_id" json:"-"`
	User       *User  `gorm:"-" json:"user"`
	Content    string `gorm:"column:content" json:"content"`
	CreateDate string `gorm:"-" json:"create_date"`
	CreateTime int64  `gorm:"column:create_time" json:"-"`
}

type commentDao struct {
}

var cDao *commentDao
var cDaoOnce sync.Once

func NewCommentOnceInstance() *commentDao {
	cDaoOnce.Do(
		func() {
			cDao = &commentDao{}
		})
	return cDao
}

func (*commentDao) GetCommentsByVideoID(vid int) []*Comment {
	comments := make([]*Comment, 0)
	db.Where("video_id=?", vid).Order("create_time DESC").Find(&comments)
	return comments
}

func (*commentDao) AddComment(com *Comment) *Comment {
	db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		video := &Video{}
		if err := tx.Create(com).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", com.VideoID).Find(video).Error; err != nil {
			return err
		}
		video.CommentCount = video.CommentCount + 1
		if err := tx.Save(video).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return com
}
