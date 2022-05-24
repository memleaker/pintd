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

const (
	CONN_DEC = iota
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
		conns int = 0
		lconn net.Conn
		rconn net.Conn
		err   error
		ch    = make(chan int)
	)

	if ch == nil {
		plog.Println("Alloc Channel Failed, listener closed.")
		return
	}

	defer wg.Done()
	defer listener.listener.Close()

	// Wait Connection coming.
	for {
		lconn, err = listener.listener.Accept()
		if err != nil {
			plog.Println("Accept Connection Failed %s, listener closed.", err.Error())
			return
		}

		// filter address
		ip, _, _ := strings.Cut(lconn.RemoteAddr().String(), ":")
		if matched := filter.FilterAddr(ip, cfg.SectionName); matched {
			lconn.Close()
			plog.Println("Matched Deny Address : %s.", ip)
			continue
		}

		// check connections number.
		for loop := true; loop; {
			select {
			case cmd := <-ch:
				if cmd == CONN_DEC {
					conns--
				}
			default:
				loop = false
			}
		}

		if conns >= cfg.MaxRedirects {
			lconn.Close()
			plog.Println("Connection Limit to %d, Closed Connection.", cfg.MaxRedirects)
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

			// Dial Success.
			conns++
			plog.Println("Dial to %s Success.", rconn.RemoteAddr().String())

			// handle data.
			go HandleData(lconn, rconn, ch)
			break
		}
	}
}
