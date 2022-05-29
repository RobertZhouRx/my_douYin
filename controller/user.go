package controller

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"myDouYIn/dao"
	"strconv"
	"time"
)

type BaseResponse struct {
	StatusCode int    `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetUserResponse struct {
	BaseResponse
	User dao.User `json:"user"`
}

type RegisterResponse struct {
	BaseResponse
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userdao := dao.NewUserOnceInstance()
	user := userdao.GetUserByUsername(username)
	if user.Username != "" {
		c.JSON(200, RegisterResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "User already exist",
			},
		})
	} else {
		salt := strconv.FormatInt(time.Now().Unix(), 10)
		hash := getHashAndSalt(password, salt)
		u := &dao.User{
			Username: username,
			Salt:     salt,
			Hash:     hash,
		}
		u = userdao.AddUserToSql(u)
		token := username + salt

		nowAccount := dao.NewNowAccountOnceInstance()
		nowAccount.UpdateAccount(token, u)
		c.JSON(200, RegisterResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
				StatusMsg:  "register success",
			},
			Token:  token,
			UserId: int64(u.Id),
		})
	}
}

func Action(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("user_id"))
	token := c.Query("token")
	account := dao.NewNowAccountOnceInstance()
	if _, ok := account.Token[token]; !ok {
		c.JSON(200, GetUserResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "token is invalid",
			},
			User: dao.User{},
		})
	} else {
		userdao := dao.NewUserOnceInstance()
		user := userdao.GetUserByUserID(uid)
		user = filterSensitive(user)
		c.JSON(200, GetUserResponse{
			BaseResponse: BaseResponse{
				StatusCode: 0,
			},
			User: *user,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	userdao := dao.NewUserOnceInstance()
	user := userdao.GetUserByUsername(username)
	if user.Username == "" {
		c.JSON(200, RegisterResponse{
			BaseResponse: BaseResponse{
				StatusCode: 1,
				StatusMsg:  "User not exist",
			},
		})
	} else {
		salt := user.Salt
		hash := user.Hash
		hashNew := getHashAndSalt(password, salt)
		if hash == hashNew {
			token := username + salt
			nowAccount := dao.NewNowAccountOnceInstance()
			nowAccount.UpdateAccount(token, user)
			c.JSON(200, RegisterResponse{
				BaseResponse: BaseResponse{
					StatusCode: 0,
					StatusMsg:  "login success",
				},
				UserId: int64(user.Id),
				Token:  token,
			})
		} else {
			c.JSON(200, RegisterResponse{
				BaseResponse: BaseResponse{
					StatusCode: 1,
					StatusMsg:  "username or password error",
				},
			})
		}
	}
}

func getHashAndSalt(password, salt string) string {
	m5 := md5.New()
	m5.Write([]byte(password))
	m5.Write([]byte(salt))
	hash := hex.EncodeToString(m5.Sum(nil))
	return hash
}

func filterSensitive(u *dao.User) *dao.User {
	u.Salt = ""
	u.Hash = ""
	return u
}
