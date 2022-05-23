package core

import (
	"net"
	"pintd/config"
	"pintd/filter"
	"pintd/plog"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Listener struct {
	listener net.Listener
}

var listeners = make(map[string]Listener, 0)

func CreateListener(cfg *config.PintdConfig) {
	for _, redirect := range cfg.Redirects {
		listener, err := net.Listen("tcp", redirect.LocalAddr+":"+strconv.Itoa(redirect.LocalPort))
		if err != nil {
			plog.Fatalln("Listen Failed : %s", err.Error())
		}

		listeners[redirect.SectionName] = Listener{listener}
	}
}

func HandleConn(cfg *config.PintdConfig) {
	var wg sync.WaitGroup

	for _, redirect := range cfg.Redirects {
		go AcceptConn(listeners[redirect.SectionName], redirect, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func AcceptConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	var (
		lconn net.Conn
		rconn net.Conn
		err   error
	)

	defer wg.Done()

	// Wait Connection coming.
	for {
		lconn, err = listener.listener.Accept()
		if err != nil {
			plog.Println("Accept Connection Failed %s, listener closed.", err.Error())
			listener.listener.Close()
			return
		}

		// filter address
		ip, _, _ := strings.Cut(lconn.RemoteAddr().String(), ":")
		if matched := filter.FilterAddr(ip, cfg.SectionName); matched {
			lconn.Close()
			plog.Println("Matched Deny Address : %s.", ip)
			continue
		}

		plog.Println("%s Accept Connection from %s.", lconn.LocalAddr().String(), lconn.RemoteAddr().String())

		// Dial to remote.
		for interval := 2; ; interval += 2 {
			if interval > 8 {
				plog.Println("Dial Failed : %s, Stop Reconnect.", err.Error())
				break
			}

			rconn, err = net.DialTimeout("tcp", cfg.RemoteAddr+":"+strconv.Itoa(cfg.RemorePort), time.Second*time.Duration(interval))
			if err != nil {
				plog.Println("Dial Failed : %s, Reconnect...", err.Error())
				continue
			}

			plog.Println("Dial to %s Success.", rconn.RemoteAddr().String())

			// handle data.
			go HandleData(lconn, rconn, cfg)
			break
		}
	}
}
