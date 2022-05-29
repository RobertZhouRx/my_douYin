package tools

import (
	"myDouYIn/dao"
	"myDouYIn/router"
	"sync"
)

func Init() {
	sy := sync.WaitGroup{}
	sy.Add(2)
	go func() {
		dao.InitDBTool()
		sy.Done()
	}()
	go func() {
		router.InitRouter()
		sy.Done()
	}()
	sy.Wait()
}
