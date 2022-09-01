package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndirectState struct {
	Id           int    `json:"id"`
	Protocol     string `json:"protocol"`
	SrcAddr      string `json:"src-addr"`
	SrcPort      string `json:"src-port"`
	DestAddr     string `json:"dest-addr"`
	DestPort     string `json:"dest-port"`
	CreateTime   string `json:"create-time"`
	ForwardFlow  string `json:"forward-flow"`
	RealTimeFlow string `json:"realtime-flow"`
}

// 这就不用数据库了
func GetIndirectState(c *gin.Context) {
	state := IndirectState{
		Id:           1,
		Protocol:     "TCP",
		SrcAddr:      "127.0.0.1",
		SrcPort:      "8888",
		DestAddr:     "127.0.0.1",
		DestPort:     "9999",
		CreateTime:   "2022-09-01",
		ForwardFlow:  "100MB",
		RealTimeFlow: "20MB/s",
	}

	state2 := IndirectState{
		Id:           1,
		Protocol:     "UDP",
		SrcAddr:      "111.111.111.111",
		SrcPort:      "8080",
		DestAddr:     "127.0.0.1",
		DestPort:     "9999",
		CreateTime:   "2022-09-01",
		ForwardFlow:  "100MB",
		RealTimeFlow: "10MB/s",
	}

	states := [10]IndirectState{state2, state, state, state2, state,
		state, state, state, state2, state2}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": 20,
		"data":  &states,
	})
}

func TerminateConn(c *gin.Context) {
	// 不但要检验id，还要看五元组是否相同，因为前端有延迟，可能删除错误
}
