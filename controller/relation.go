package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"strconv"
)

type RelationActionResponse struct {
	BaseResponse
}

type RelationFollowListResponse struct {
	BaseResponse
	UserList []*dao.User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId, _ := strconv.Atoi(c.Query("to_user_id"))
	actionType := c.Query("action_type")
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
	fDao := dao.NewFollowerOnceInstance()
	follow := &dao.Follower{
		FollowID:   toUserId,
		FollowerID: uid,
	}
	if actionType == "1" {
		err := fDao.AddFollower(follow)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(200, PublishActionResponse{
				BaseResponse: BaseResponse{
					StatusCode: 1,
				}})
		}
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
			}})
	} else {
		err := fDao.DeleteFollower(follow)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(200, PublishActionResponse{
				BaseResponse: BaseResponse{
					StatusCode: 1,
				}})
		}
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
			}})
	}
}

func RelationFollowList(c *gin.Context) {
	token := c.Query("token")
	account := dao.NewNowAccountOnceInstance()
	uid, _ := strconv.Atoi(c.Query("user_id"))
	if _, ok := account.Token[token]; !ok {
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			}})
		return
	}
	fDao := dao.NewFollowerOnceInstance()
	relations := fDao.GetCommentsByFollowerId(uid)
	users := getUserByFollowID(relations, uid)
	c.JSON(200, RelationFollowListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
		},
		UserList: users,
	},
	)
}
func getUserByFollowID(rels []*dao.Follower, uid int) []*dao.User {
	uDao := dao.NewUserOnceInstance()
	var u *dao.User
	users := make([]*dao.User, 0)
	for i := 0; i < len(rels); i++ {
		u = uDao.GetUserByUserID(rels[i].FollowID)
		SetIsFollower(rels[i].FollowID, uid, u)
		users = append(users, u)
	}
	return users
}

func RelationFollowerList(c *gin.Context) {
	token := c.Query("token")
	account := dao.NewNowAccountOnceInstance()
	uid, _ := strconv.Atoi(c.Query("user_id"))
	if _, ok := account.Token[token]; !ok {
		c.JSON(200, PublishActionResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			}})
		return
	}
	fDao := dao.NewFollowerOnceInstance()
	relations := fDao.GetCommentsByFollowId(uid)
	users := getUserByFollowerID(relations, uid)
	c.JSON(200, RelationFollowListResponse{
		BaseResponse: BaseResponse{
			StatusCode: 0,
		},
		UserList: users,
	},
	)
}
func getUserByFollowerID(rels []*dao.Follower, uid int) []*dao.User {
	uDao := dao.NewUserOnceInstance()
	var u *dao.User
	users := make([]*dao.User, 0)
	for i := 0; i < len(rels); i++ {
		u = uDao.GetUserByUserID(rels[i].FollowerID)
		SetIsFollower(rels[i].FollowerID, uid, u)
		users = append(users, u)
	}
	return users
}
