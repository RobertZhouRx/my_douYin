package dao

import "sync"

type Video struct {
	Id            int    `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Title         string `gorm:"column:title" json:"title"`
	PlayUrl       string `gorm:"column:play_url" json:"play_url"`
	AuthorId      int    `gorm:"column:author_id" json:"author_id"`
	CoverUrl      string `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount int    `gorm:"column:favorite_count" json:"favorite_count"`
	CommentCount  int    `gorm:"column:comment_count" json:"comment_count"`
	CreateTime    int64  `gorm:"column:create_time" json:"create_time"`
	IsFavorite    bool   `gorm:"-"  json:"is_favorite"`
	User          *User  `gorm:"-" json:"author"`
}

type videoDao struct {
}

var vDao *videoDao
var videoOnce sync.Once

func NewVideoOnceInstance() *videoDao {
	videoOnce.Do(
		func() {
			vDao = &videoDao{}
		})
	return vDao
}

func (*videoDao) AddVideoToSql(video *Video) *Video {
	db.Create(video)
	return video
}

func (*videoDao) GetVideosByUid(uid int) []*Video {
	videos := make([]*Video, 0)
	db.Where("author_id = ?", uid).Find(&videos)
	return videos
}

func (*videoDao) GetVideoByVid(vid int) *Video {
	video := &Video{}
	db.Where("id = ?", vid).Find(&video)
	return video
}

func (*videoDao) UpdateVideo(video *Video) {
	db.Save(video)
}

func (*videoDao) GetVideos(lastTime int64) []*Video {
	videos := make([]*Video, 0)
	db.Where("create_time<?", lastTime).Limit(10).Order("create_time DESC").Find(&videos)
	return videos
}
