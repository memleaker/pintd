package main

import (
	"pintd/config"
	"pintd/core"
	"pintd/filter"
	"pintd/plog"
)

func main() {
	// read config.
	cfg := config.ReadConfig()

	// init log module.
	plog.InitLog(cfg.LogFile)

	// init deny rules.
	filter.AddDenyAddrs(cfg)

	// create listener.
	core.CreateListener(cfg)

	// listen and running...
	core.HandleConn(cfg)
}
