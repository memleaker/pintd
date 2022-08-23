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
