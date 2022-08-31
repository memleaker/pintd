package api

import (
	"net/http"
	"pintd/model"
	"pintd/plog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLog(c *gin.Context) {
	// get url args.
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page <= 0 || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误的请求参数",
		})
		return
	}

	// get rows number.
	rows, err := model.GetLogTblRows()
	if err != nil {
		plog.Println("model.GetLogTblRows() : " + err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})
		return
	}

	logs, err := model.GetLog(page, limit)
	if err != nil {
		// get config failed.
		plog.Println("model.GetLog() : " + err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": 0,
			"msg":     "服务器发生错误!",
		})
		return
	}

	// get config success.
	c.JSON(http.StatusOK, gin.H{
		"code":  0, // code:0, means no error.
		"count": rows,
		"data":  &logs,
	})
}

func DelLog(c *gin.Context) {
	log := model.Logging{}

	// parse JSON to structure.
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	// delete config.
	ok, err := model.DelLog(log.Id)
	if !ok {
		// delete failed.
		if err != nil {
			plog.Println("model.DelLog() : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "服务器发生错误!",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "删除配置失败",
			})
		}

		return
	}

	// delete success.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "删除配置成功",
	})
}

func DelMoreLog(c *gin.Context) {
	// get url args.
	num, _ := strconv.Atoi(c.Query("num"))
	if num <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误的请求参数",
		})
		return
	}

	logs := make([]model.Logging, num)

	// parse JSON to structure.
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": 0,
			"msg":     "错误请求格式",
		})
		return
	}

	// delete config.
	ok, err := model.DelMoreLog(logs)
	if !ok {
		// delete failed.
		if err != nil {
			plog.Println("model.DelMoreLog() : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "服务器发生错误!",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": 0,
				"msg":     "删除配置失败",
			})
		}

		return
	}

	// delete success.
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"msg":     "删除配置成功",
	})
}
