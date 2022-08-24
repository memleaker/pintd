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
	fmt.Println(cfg)

	// return.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "新建配置成功",
	})
}

func IndirectCfgShow(c *gin.Context) {
	cfg := [1]model.IndirectConfig{{
		Protocol:   "TCP",
		ListenAddr: "127.0.0.1",
		ListenPort: "8888",
		DestAddr:   "127.0.0.1",
		DestPort:   "9999",
		Acl:        "黑名单",
		AdmitAddr:  "",
		DenyAddr:   "1.1.1.1",
		MaxConns:   "100",
		Memo:       "测试测试"}}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": &cfg,
	})
}
