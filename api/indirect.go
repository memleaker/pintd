package api

import (
	"net/http"
	"pintd/model"
	"strconv"

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
			"msg":     "错误请求格式",
		})
		return
	}

	// check values. 检查长度, IP, 端口
	if model.RepeatConfig(&cfg) {
		c.JSON(http.StatusConflict, gin.H{
			"success": 0,
			"msg":     "配置已经存在",
		})
		return
	}

	// new config.
	if model.NewIndirectCfg(&cfg) {
		c.JSON(http.StatusCreated, gin.H{
			"success": 1,
			"msg":     "新建配置成功",
		})
		return
	}

	// create failed.
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": 0,
		"msg":     "新建配置失败",
	})
}

func IndirectCfgDel(c *gin.Context) {
	// 删除配置只需要校验协议和监听的端口即可, 删除运行状态的具体连接才需要校验五元组
}

func IndirectCfgShow(c *gin.Context) {

	// parse url args.
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page <= 0 || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误的请求参数",
		})
		return
	}

	cfgs := model.GetIndirectCfg(page, limit)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,         // code:0 表示没有错误
		"count": len(cfgs), // 表示共有多少条
		"data":  &cfgs,     // 此处应当解析客户端limit参数, 返回小于等于limit参数个条目
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
