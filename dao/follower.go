package dao

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type Follower struct {
	ID         int `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"` // 主键id
	FollowID   int `gorm:"column:follow_id" json:"follow_id"`
	FollowerID int `gorm:"column:follower_id" json:"follower_id"`
}

type followerDao struct {
}

var fDao *followerDao
var fDaoOnce sync.Once

func NewFollowerOnceInstance() *followerDao {
	userOnce.Do(
		func() {
			fDao = &followerDao{}
		})
	return fDao
}

func (*followerDao) GetCommentsByFollowerId(uid int) []*Follower {
	followers := make([]*Follower, 0)
	db.Where("follower_id=?", uid).Find(&followers)
	return followers
}

func (*followerDao) GetCommentsByFollowId(uid int) []*Follower {
	followers := make([]*Follower, 0)
	db.Where("follow_id=?", uid).Find(&followers)
	return followers
}

func (*followerDao) AddFollower(f *Follower) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		follow := &User{}
		follower := &User{}
		if err := tx.Create(f).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.FollowID).Find(follow).Error; err != nil {
			return err
		}
		follow.FollowerCount = follow.FollowerCount + 1
		if err := tx.Save(follow).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.FollowerID).Find(follower).Error; err != nil {
			return err
		}
		follower.FollowCount = follower.FollowCount + 1
		if err := tx.Save(follower).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (*followerDao) DeleteFollower(f *Follower) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		follow := &User{}
		follower := &User{}
		res := tx.Delete(&Follower{}, "follow_id=? and follower_id=?", f.FollowID, f.FollowerID)
		if res.Error != nil || res.RowsAffected == 0 {
			return errors.New("delete failure")
		}
		if err := tx.Where("id=?", f.FollowID).Find(follow).Error; err != nil {
			return err
		}
		follow.FollowerCount = follow.FollowerCount - 1
		if err := tx.Save(follow).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.FollowerID).Find(follower).Error; err != nil {
			return err
		}
		follower.FollowCount = follower.FollowCount - 1
		if err := tx.Save(follower).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (*followerDao) GetFollowerByIDs(followID, followerId int) *Follower {
	follow := &Follower{}
	db.Where("follow_id=? and follower_id=?", followID, followerId).Find(follow)
	return follow
}
