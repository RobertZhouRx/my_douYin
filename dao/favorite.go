package dao

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type Favorite struct {
	UserID  int `gorm:"column:user_id;primary_key;" json:"user_id"`
	VideoID int `gorm:"column:video_id;primary_key;" json:"video_id"`
}

type favoriteDao struct {
}

var faDao *favoriteDao
var faDaoOnce sync.Once

func NewFavoriteOnceInstance() *favoriteDao {
	userOnce.Do(
		func() {
			faDao = &favoriteDao{}
		})
	return faDao
}

func (*favoriteDao) GetFavorite(userID, videoID int) *Favorite {
	favorite := &Favorite{}
	db.Where("user_id=? and video_id=?", userID, videoID).Find(favorite)
	return favorite
}

func (*favoriteDao) GetFavoritesByUid(userID int) []*Favorite {
	favorites := make([]*Favorite, 0)
	db.Where("user_id=?", userID).Find(&favorites)
	return favorites
}

func (*favoriteDao) AddFavoriteToSql(f *Favorite) {
	db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		user := &User{}
		video := &Video{}
		if err := tx.Create(f).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.UserID).Find(user).Error; err != nil {
			return err
		}
		user.FavoriteCount = user.FavoriteCount + 1
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.VideoID).Find(video).Error; err != nil {
			return err
		}
		video.FavoriteCount = video.FavoriteCount + 1
		if err := tx.Save(video).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (*favoriteDao) DeleteFavoriteFromSql(f *Favorite) {
	db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		user := &User{}
		video := &Video{}
		res := tx.Delete(f)
		if res.Error != nil || res.RowsAffected == 0 {
			return errors.New("delete failure")
		}

		if err := tx.Where("id=?", f.UserID).Find(user).Error; err != nil {
			return err
		}
		user.FavoriteCount -= 1
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		if err := tx.Where("id=?", f.VideoID).Find(video).Error; err != nil {
			return err
		}
		video.FavoriteCount -= 1
		if err := tx.Save(video).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}
