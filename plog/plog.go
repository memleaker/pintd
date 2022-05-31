package plog

import (
	"fmt"
	"log"
	"os"
	"pintd/config"
)

var logfile *os.File
var plog *log.Logger
var appmode int

func InitLog(cfg *config.PintdConfig) {
	var err error

	logfile, err = os.OpenFile(cfg.LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("Open Log File Failed : ", err.Error())
	}

	plog = log.New(logfile, "", log.Ldate|log.Ltime)

	if cfg.AppMode == "release" {
		appmode = 1
	} else {
		appmode = 0
	}
}

func Println(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	plog.Output(2, s)

	if appmode == 0 {
		log.Println(s)
	}
}

func Fatalln(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	plog.Output(2, s)
	log.Fatalln(s)
}
