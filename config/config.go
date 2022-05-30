package config

import (
	"log"
	"net"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	CONFIGFILE  = "pintd.ini"
	DEBUGMODE   = "debug"
	RELEASEMODE = "release"
	LOGFILE     = "/var/log/pintd.log"
)

type PintdConfig struct {
	AppMode   string
	LogFile   string
	Redirects []*RedirectConfig
}

type RedirectConfig struct {
	LocalPort    int
	RemorePort   int
	MaxRedirects int
	Protocol     string
	LocalAddr    string
	RemoteAddr   string
	SectionName  string
	Denyaddr     []string
}

// read pintd and indirect config.
func ReadConfig(cfgfile string) *PintdConfig {
	var Pintdcf PintdConfig

	Pintdcf.Redirects = make([]*RedirectConfig, 0)
	if Pintdcf.Redirects == nil {
		log.Fatalln("Alloc structure []RedirectConfig Failed.")
	}

	cf, err := ini.Load(cfgfile)
	if err != nil {
		log.Fatalln("Read Config Failed :", err.Error())
	}

	log.Println("Using Config File", cfgfile)

	// read pintd config.
	Pintdcf.AppMode = cf.Section("pintd").Key("AppMode").In("debug", []string{"debug", "release"})
	Pintdcf.LogFile = cf.Section("pintd").Key("LogFile").MustString(LOGFILE)
	if Pintdcf.AppMode == DEBUGMODE {
		log.Println("Pintd is Running On debug mode.")
	}

	// read port redirect config.
	childs := cf.Section("redirect").ChildSections()

	for _, section := range childs {
		redirect := new(RedirectConfig)
		if redirect == nil {
			log.Fatalln("Alloc structure RedirectConfig Failed.")
		}

		redirect.LocalAddr = section.Key("localaddr").MustString("0.0.0.0")
		if addr := net.ParseIP(redirect.LocalAddr); addr == nil {
			log.Fatalln("Invalid Addr : ", redirect.LocalAddr)
		}

		redirect.LocalPort = section.Key("localport").MustInt(8888)
		if redirect.LocalPort < 0 || redirect.LocalPort > 65535 {
			log.Fatalln("Invalid Port : ", redirect.LocalPort)
		}

		redirect.RemoteAddr = section.Key("remoteaddr").MustString("127.0.0.1")
		if addr := net.ParseIP(redirect.RemoteAddr); addr == nil {
			log.Fatalln("Invalid Addr : ", redirect.RemoteAddr)
		}

		redirect.RemorePort = section.Key("remoteport").MustInt(9999)
		if redirect.RemorePort < 0 || redirect.RemorePort > 65535 {
			log.Fatalln("Invalid Port : ", redirect.RemorePort)
		}

		redirect.SectionName = section.Name()

		redirect.Protocol = section.Key("proto").MustString("tcp")
		if redirect.Protocol != "tcp" && redirect.Protocol != "udp" {
			log.Fatalln("Invalid Protocol : ", redirect.Protocol)
		}

		redirect.MaxRedirects = section.Key("maxredirects").MustInt(100)

		redirect.Denyaddr = section.Key("denyaddrs").Strings(",")
		for _, addr := range redirect.Denyaddr {
			before, after, ok := strings.Cut(addr, "/")
			if ok {
				addr = before
				mask, _ := strconv.Atoi(after)
				if mask <= 0 || mask > 32 {
					log.Fatalln("Invalid Mask : ", after)
				}
			}

			if ip := net.ParseIP(addr); ip == nil {
				log.Fatalln("Invalid Addr : ", addr)
			}
		}

		Pintdcf.Redirects = append(Pintdcf.Redirects, redirect)
	}

	return &Pintdcf
}
