package router

import (
	"pintd/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode("debug")

	engine := gin.Default()

	engine.LoadHTMLFiles("./web/index.html")
	engine.Static("/layui", "./web/layui")
	engine.Static("/static", "./web/static")

	// router group
	router := engine.Group("/")
	{
		// GET
		router.GET("", api.MainPage)

		// POST
		router.POST("/indirect/cfg_new", api.IndirectCfgNew)
	}

	return engine
}
