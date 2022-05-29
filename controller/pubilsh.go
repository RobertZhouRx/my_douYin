package controller

import (
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"path"
	"strconv"
	"time"
)

type PublishActionResponse struct {
	BaseResponse
}

type PublishListResponse struct {
	BaseResponse
	VideoList []*dao.Video `json:"video_list"`
}

func PublishList(c *gin.Context) {
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
	videoDao := dao.NewVideoOnceInstance()
	videos := videoDao.GetVideosByUid(uid)
	c.JSON(200, PublishListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}

func PublishLogin(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	account := dao.NewNowAccountOnceInstance()
	uid, ok := account.Token[token]
	if !ok {
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			}})
		return
	}
	videoDao := dao.NewVideoOnceInstance()
	video := videoDao.AddVideoToSql(&dao.Video{
		Title:      title,
		AuthorId:   uid,
		CreateTime: time.Now().Unix(),
	})
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "video is invalid",
			}})
		return
	}

	url := "http://192.168.0.101:2222/" + "douyin/video/" + "?palyURL="
	videoUrl := "public/video" + strconv.FormatInt(int64(video.Id), 10) + path.Ext(file.Filename)
	c.SaveUploadedFile(file, videoUrl)
	video.PlayUrl = url + videoUrl
	video.CoverUrl = url + "public/1.jpeg"
	videoDao.UpdateVideo(video)

	c.JSON(200, PublishActionResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "upload success",
		}})
}
