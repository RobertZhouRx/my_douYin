package controller

import (
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"strconv"
	"time"
)

type CommentActionResponse struct {
	BaseResponse
	Comment *dao.Comment `json:"comment"`
}

type CommentListResponse struct {
	BaseResponse
	CommentList []*dao.Comment `json:"comment_list"`
}

func CommentList(c *gin.Context) {
	vid, _ := strconv.Atoi(c.Query("video_id"))
	cDao := dao.NewCommentOnceInstance()
	comments := cDao.GetCommentsByVideoID(vid)
	fillComments(comments)
	c.JSON(200, CommentListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CommentList: comments,
	})
}

func fillComments(com []*dao.Comment) {
	uDao := dao.NewUserOnceInstance()
	for i := 0; i < len(com); i++ {
		t := time.Unix(com[i].CreateTime, 0)
		com[i].User = uDao.GetUserByUserID(com[i].UserID)
		com[i].CreateDate = t.Format("01-02")
	}
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	vid, _ := strconv.Atoi(c.Query("video_id"))
	text := c.Query("comment_text")
	var uid int
	var ok bool
	account := dao.NewNowAccountOnceInstance()
	if uid, ok = account.Token[token]; !ok {
		c.JSON(200, CommentActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			}})
		return
	}
	cDao := dao.NewCommentOnceInstance()
	comment := &dao.Comment{
		UserID:     uid,
		VideoID:    vid,
		Content:    text,
		CreateTime: time.Now().Unix(),
	}
	comment = cDao.AddComment(comment)
	c.JSON(200, CommentActionResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
			StatusMsg:  "token is invalid",
		},
		Comment: comment,
	},
	)
}
