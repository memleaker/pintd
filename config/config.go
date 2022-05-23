package config

import (
	"log"

	"gopkg.in/ini.v1"
)

const (
	CONFIGFILE  = "pintd.ini"
	DEBUGMODE   = "debug"
	RELEASEMODE = "release"
	LOGFILE     = "/var/log/pintd.log"
)

type PintdConfig struct {
	AppMode      string
	MaxRedirects int
	LogFile      string
	Redirects    []*RedirectConfig
}

type RedirectConfig struct {
	LocalPort   int
	RemorePort  int
	LocalAddr   string
	RemoteAddr  string
	SectionName string
	Denyaddr    []string
}

// read pintd and indirect config.
func ReadConfig() *PintdConfig {
	var Pintdcf PintdConfig

	Pintdcf.Redirects = make([]*RedirectConfig, 0)
	if Pintdcf.Redirects == nil {
		log.Fatalln("Alloc structure []RedirectConfig Failed.")
	}

	cf, err := ini.Load(CONFIGFILE)
	if err != nil {
		log.Fatalln("Read Config Failed :", err.Error())
	}

	// read pintd config.
	Pintdcf.AppMode = cf.Section("pintd").Key("AppMode").In("debug", []string{"debug", "release"})
	Pintdcf.MaxRedirects = cf.Section("pintd").Key("MaxRedirects").MustInt(1024)
	Pintdcf.LogFile = cf.Section("pintd").Key("LogFile").MustString(LOGFILE)
	if Pintdcf.AppMode == DEBUGMODE {
		log.Println("Pintd is Running On debug mode.")
	}

	// read port redirect config.
	childs := cf.Section("redirect").ChildSections()

	for index, section := range childs {
		if index > Pintdcf.MaxRedirects {
			log.Println("Too much redirect condig, some will be skipped.")
			break
		}

		redirect := new(RedirectConfig)
		if redirect == nil {
			log.Fatalln("Alloc structure RedirectConfig Failed.")
		}

		redirect.LocalAddr = section.Key("localaddr").MustString("0.0.0.0")
		redirect.LocalPort = section.Key("localport").MustInt(8888)
		redirect.RemoteAddr = section.Key("remoteaddr").MustString("127.0.0.1")
		redirect.RemorePort = section.Key("remoteport").MustInt(9999)
		redirect.SectionName = section.Name()
		redirect.Denyaddr = section.Key("denyaddrs").Strings(",")

		Pintdcf.Redirects = append(Pintdcf.Redirects, redirect)
	}

	// todo CheckConfig

	return &Pintdcf
}
