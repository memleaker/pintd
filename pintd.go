package main

import (
	"flag"
	"pintd/config"
	"pintd/core"
	"pintd/filter"
	"pintd/plog"
)

func main() {

	cfgfile := flag.String("c", config.CONFIGFILE, "use -c to specific [config file]")

	flag.Parse()

	// read config.
	cfg := config.ReadConfig(*cfgfile)

	// init log module.
	plog.InitLog(cfg.LogFile)

	// init deny rules.
	filter.AddDenyAddrs(cfg)

	// create listener.
	core.CreateListener(cfg)

	// listen and running...
	core.HandleConns(cfg)
}
