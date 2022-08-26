package api

import (
	"fmt"
	"net/http"
	"pintd/model"

	"github.com/gin-gonic/gin"
)

// Indirect Config.

// Indirect New Config.
func IndirectCfgNew(c *gin.Context) {
	cfg := model.IndirectConfig{}

	// parse JSON.
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求",
		})

		return
	}

	// check values.
	// 检查输入，消息是否过长
	// 只需检查协议和本地监听端口是否重复即可, 其他不必
	fmt.Println(cfg)

	// return.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "新建配置成功",
	})
}

func IndirectCfgDel(c *gin.Context) {
	// 删除配置只需要校验协议和监听的端口即可, 删除运行状态的具体连接才需要校验五元组
}

func IndirectCfgShow(c *gin.Context) {
	cfg := model.IndirectConfig{
		Protocol:   "TCP",
		ListenAddr: "127.0.0.1",
		ListenPort: "8888",
		DestAddr:   "127.0.0.1",
		DestPort:   "9999",
		Acl:        "黑名单",
		AdmitAddr:  "",
		DenyAddr:   "1.1.1.1",
		MaxConns:   "100",
		Memo:       "测试测试"}

	// page1
	cfgs := [10]model.IndirectConfig{cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg}

	// page2
	cfgs_pg2 := [10]model.IndirectConfig{cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg, cfg}

	// todo 应该解析page参数,返回给客户端第几页

	c.JSON(http.StatusOK, gin.H{
		"code":  0,                         // code:0 表示没有错误
		"count": len(cfgs) + len(cfgs_pg2), // 表示共有多少条
		"data":  &cfgs,                     // 此处应当解析客户端limit参数, 返回小于等于limit参数个条目
	})
}

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
