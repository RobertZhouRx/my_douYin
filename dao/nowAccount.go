package dao

import (
	"sync"
)

type nowAccount struct {
	NowUser map[int]*User
	Token   map[string]int
}

var account *nowAccount
var accountOnce sync.Once

func NewNowAccountOnceInstance() *nowAccount {
	accountOnce.Do(
		func() {
			account = &nowAccount{
				make(map[int]*User, 0),
				make(map[string]int, 0),
			}
			account.Token["aaaaaa1653464011"] = 12
		})
	return account
}

func (acc *nowAccount) UpdateAccount(token string, user *User) {
	acc.NowUser[user.Id] = user
	acc.Token[token] = user.Id
}
