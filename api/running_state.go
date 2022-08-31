package api

import (
	"net/http"
	"pintd/model"

	"github.com/gin-gonic/gin"
)

// 这就不用数据库了
func IndirectState(c *gin.Context) {
	state := model.IndirectState{
		Protocol:     "TCP",
		SrcAddr:      "127.0.0.1",
		SrcPort:      "8888",
		DestAddr:     "127.0.0.1",
		DestPort:     "9999",
		RunningTime:  "10min",
		ForwardFlow:  "100MB",
		RealTimeFlow: "20MB/s",
	}

	state2 := model.IndirectState{
		Protocol:     "UDP",
		SrcAddr:      "111.111.111.111",
		SrcPort:      "8080",
		DestAddr:     "127.0.0.1",
		DestPort:     "9999",
		RunningTime:  "10min",
		ForwardFlow:  "100MB",
		RealTimeFlow: "10MB/s",
	}

	states := [10]model.IndirectState{state2, state, state, state2, state,
		state, state, state, state2, state2}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,       // code:0 表示没有错误
		"count": 20,      // 表示共有多少条
		"data":  &states, // 此处应当解析客户端limit参数, 返回小于等于limit参数个条目
	})
}
