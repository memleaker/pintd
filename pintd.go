package main

import (
	"flag"
	"log"
	"pintd/config"
	"pintd/core"
	"pintd/filter"
	"pintd/plog"
	"syscall"
)

func main() {

	// arguments
	cfgfile := flag.String("c", config.CONFIGFILE, "use -c to specific [config file]")
	flag.Parse()

	// read config.
	cfg := config.ReadConfig(*cfgfile)

	// init system, ulmits
	// golang don't need ignore SIGPIPE, but write return EPIPE
	SetSystem(cfg)

	// init log module.
	plog.InitLog(cfg)

	// init rules.
	filter.InitFilter(cfg)

	// create listeners.
	listeners := core.InitListeners(cfg)

	// listen and running...
	core.HandleConns(cfg, listeners)
}

func SetSystem(cfg *config.PintdConfig) {

	// set open fd numbers
	rlim := syscall.Rlimit{Cur: cfg.MaxOpenFiles, Max: cfg.MaxOpenFiles}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlim); err != nil {
		log.Fatalln("Set resources limit failed : " + err.Error())
	}
}
