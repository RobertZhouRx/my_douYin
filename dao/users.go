package dao

import (
	"sync"
)

type User struct {
	Id            int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 主键id
	Name          string `gorm:"column:name;NOT NULL" json:"-"`                  // 用户昵称
	Avatar        string `gorm:"column:avatar;NOT NULL" json:"avatar"`           // 头像
	Username      string `gorm:"column:username;unique;NOT NULL" json:"name"`
	Hash          string `gorm:"column:hash;NOT NULL" json:"hash"`
	Salt          string `gorm:"column:salt;NOT NULL" json:"salt"`
	VideoCount    int    `gorm:"column:video_count;NOT NULL" json:"video_count"`
	FollowCount   int    `gorm:"column:follow_count;NOT NULL" json:"follow_count"`
	FollowerCount int    `gorm:"column:follower_count;NOT NULL" json:"follower_count"`
	FavoriteCount int    `gorm:"column:favorite_count;NOT NULL" json:"favorite_count"`
	IsFollow      bool   `gorm:"-" json:"is_follow"`
}

type userDao struct {
}

var uDao *userDao
var userOnce sync.Once

func NewUserOnceInstance() *userDao {
	userOnce.Do(
		func() {
			uDao = &userDao{}
		})
	return uDao
}

func (*userDao) AddUserToSql(user *User) *User {
	db.Create(user)
	return user
}

func (*userDao) GetUserByUsername(username string) *User {
	user := &User{}
	db.Where("username=?", username).Find(user)
	return user
}

func (*userDao) GetUserByUserID(id int) *User {
	user := &User{}
	db.Where("id=?", id).Find(user)
	return user
}
