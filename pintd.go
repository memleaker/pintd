package main

import (
	"flag"
	"pintd/config"
	"pintd/model"
	"pintd/plog"
	"pintd/router"
)

func main() {

	cfgfile := flag.String("c", config.CONFIGFILE, "use -c to specific [config file]")

	flag.Parse()

	// read config.
	cfg := config.ReadConfig(*cfgfile)

	// init log module.
	plog.InitLog(cfg)

	// init deny rules.
	//filter.AddDenyAddrs(cfg)

	// create listener.
	//core.CreateListener(cfg)

	// listen and running...
	//core.HandleConns(cfg)

	// initialize database.
	model.InitDb("/root/Go/pintd/pintd.db")

	r := router.InitRouter()

	r.Run(":8888")
}
