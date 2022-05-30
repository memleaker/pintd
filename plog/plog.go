package plog

import (
	"fmt"
	"log"
	"os"
)

var logfile *os.File
var plog *log.Logger

func InitLog(file string) {
	var err error

	logfile, err = os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("Open Log File Failed : ", err.Error())
	}

	plog = log.New(logfile, "", log.Ldate|log.Ltime)
}

func Println(format string, a ...any) {
	plog.Output(2, fmt.Sprintf(format, a...))
}

func Fatalln(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	plog.Output(2, s)
	log.Fatalln(s)
}
