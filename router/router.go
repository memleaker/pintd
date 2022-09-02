package router

import (
	"pintd/api"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode("debug")

	engine := gin.Default()

	engine.LoadHTMLFiles("./web/index.html")
	engine.Static("/web/layui", "./web/layui")
	engine.Static("/web/static", "./web/static")

	// router group
	router := engine.Group("/")
	{
		// Index.
		router.GET("", api.MainPage)

		// running state.
		router.GET("/running/state", api.GetRunningState)
		router.GET("/running/runninginfo", api.GetRunningInfo)
		router.POST("/running/terminate", api.TerminateConn)

		// indirect.
		router.GET("/indirect/cfg_show", api.IndirectCfgShow)
		router.POST("/indirect/cfg_new", api.IndirectCfgNew)
		router.POST("/indirect/cfg_del", api.IndirectCfgDel)
		router.POST("/indirect/cfg_edit", api.IndirectCfgEdit)

		// log.
		router.GET("/logging/get", api.GetLog)
		router.POST("/logging/del", api.DelLog)
		router.POST("/logging/delmore", api.DelMoreLog)
	}

	return engine
}
