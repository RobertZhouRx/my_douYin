package controller

import (
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"strconv"
)

type FavoriteActionResponse struct {
	BaseResponse
}

type FavoriteListResponse struct {
	BaseResponse
	VideoList []*dao.Video `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	vid, _ := strconv.Atoi(c.Query("video_id"))
	actionType := c.Query("action_type")
	account := dao.NewNowAccountOnceInstance()
	faDao := dao.NewFavoriteOnceInstance()
	uid, ok := account.Token[token]
	if !ok {
		c.JSON(200, FavoriteActionResponse{BaseResponse{
			StatusCode: 1,
			StatusMsg:  "token is invalid",
		}})
		return
	}
	f := &dao.Favorite{
		UserID:  uid,
		VideoID: vid,
	}
	if actionType == "1" {
		faDao.AddFavoriteToSql(f)
	} else {
		faDao.DeleteFavoriteFromSql(f)
	}
	c.JSON(200, FavoriteActionResponse{BaseResponse{
		StatusCode: 0,
	}})
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	uid, _ := strconv.Atoi(c.Query("user_id"))
	account := dao.NewNowAccountOnceInstance()
	if _, ok := account.Token[token]; !ok {
		c.JSON(200, PublishListResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			}})
		return
	}
	favoriteDao := dao.NewFavoriteOnceInstance()
	favorites := favoriteDao.GetFavoritesByUid(uid)
	videos := getVideosFromFavorites(favorites)
	c.JSON(200, PublishListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}

func getVideosFromFavorites(favorites []*dao.Favorite) []*dao.Video {
	videos := make([]*dao.Video, 0)
	videoDao := dao.NewVideoOnceInstance()
	for i := 0; i < len(favorites); i++ {
		video := videoDao.GetVideoByVid(favorites[i].VideoID)
		videos = append(videos, video)
	}
	return videos
}
