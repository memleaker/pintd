package core

import (
	"net"
	"pintd/config"
	"pintd/plog"
	"strconv"
	"sync"
)

type Listener struct {
	listener net.Listener // for tcp
	udpconn  *net.UDPConn // for udp
}

var listeners = make(map[string]Listener, 0)

func CreateListener(cfg *config.PintdConfig) {
	var (
		err      error
		listener net.Listener
		udpconn  *net.UDPConn
	)

	for _, redirect := range cfg.Redirects {
		if redirect.Protocol == "tcp" {
			listener, err = net.Listen("tcp", redirect.LocalAddr+":"+strconv.Itoa(redirect.LocalPort))
			if err != nil {
				plog.Fatalln("Listen Failed : %s", err.Error())
			}

			listeners[redirect.SectionName] = Listener{listener, nil}
		} else if redirect.Protocol == "udp" {
			udpconn, err = net.ListenUDP("udp", &net.UDPAddr{
				IP:   net.ParseIP(redirect.LocalAddr),
				Port: redirect.LocalPort})
			if err != nil {
				plog.Fatalln("ListenUDP Failed : %s", err.Error())
			}

			listeners[redirect.SectionName] = Listener{nil, udpconn}
		}
	}
}

func HandleConns(cfg *config.PintdConfig) {
	var wg sync.WaitGroup

	for _, redirect := range cfg.Redirects {

		if redirect.Protocol == "tcp" {
			go HandleTcpConn(listeners[redirect.SectionName], redirect, &wg)
			wg.Add(1)
		} else if redirect.Protocol == "udp" {
			go HandleUdpConn(listeners[redirect.SectionName], redirect, &wg)
			wg.Add(1)
		}
	}

	wg.Wait()
}
