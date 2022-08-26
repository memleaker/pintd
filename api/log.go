package api

import (
	"net/http"
	"pintd/model"

	"github.com/gin-gonic/gin"
)

func GetLog(c *gin.Context) {
	log := [1]model.Logging{{Id: "1", Time: "2022-08-26:17:32:56", Content: "test"}}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"count": 1,
		"data":  &log,
	})
}

func DelLog(c *gin.Context) {

}
