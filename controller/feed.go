package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"strconv"
)

type FeedActionResponse struct {
	BaseResponse
	NextTime  int64        `json:"next_time"`
	VideoList []*dao.Video `json:"video_list"`
}

func FeedAction(c *gin.Context) {
	lateTime := c.Query("latest_time")
	token := c.Query("token")
	late, err := strconv.ParseInt(lateTime, 10, 64)
	account := dao.NewNowAccountOnceInstance()
	uid := account.Token[token]
	if err != nil {
		fmt.Println(err)
	}
	videoDao := dao.NewVideoOnceInstance()
	videos := videoDao.GetVideos(late)
	if len(videos) == 0 {
		c.JSON(200,
			FeedActionResponse{
				BaseResponse: BaseResponse{
					StatusCode: 1,
					StatusMsg:  "video list is nil",
				},
			})
		return
	}
	SetFavoriteByVideos(videos, uid)
	setUer(uid, videos)
	c.JSON(200,
		FeedActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
			},
			VideoList: videos,
			NextTime:  videos[len(videos)-1].CreateTime,
		})
}

func FeedVideo(c *gin.Context) {
	path := c.Query("palyURL")
	c.File(path)
}

func SetFavoriteByVideos(videos []*dao.Video, uid int) {
	favoriteDao := dao.NewFavoriteOnceInstance()
	var f *dao.Favorite
	for i := 0; i < len(videos); i++ {
		f = favoriteDao.GetFavorite(uid, videos[i].Id)
		if f.UserID != 0 {
			videos[i].IsFavorite = true
		}
	}
}

func SetFavoriteByVideo(videos *dao.Video, uid int) {
	favoriteDao := dao.NewFavoriteOnceInstance()
	f := favoriteDao.GetFavorite(uid, videos.Id)
	if f.UserID != 0 {
		videos.IsFavorite = true
	}
}

func SetIsFollower(followId, followerId int, user *dao.User) {
	followDao := dao.NewFollowerOnceInstance()
	rel := followDao.GetFollowerByIDs(followId, followerId)
	if rel.ID > 0 {
		user.IsFollow = true
	}
}

func setUer(uid int, videos []*dao.Video) {
	uDao := dao.NewUserOnceInstance()
	var u *dao.User
	for i := 0; i < len(videos); i++ {
		u = uDao.GetUserByUserID(videos[i].AuthorId)
		SetIsFollower(videos[i].AuthorId, uid, u)
		videos[i].User = u
	}
}
