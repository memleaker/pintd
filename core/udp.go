package core

import (
	"errors"
	"net"
	"os"
	"pintd/config"
	"pintd/plog"
	"sync"
	"time"
)

var conns map[string]chan *Dgram = make(map[string]chan *Dgram)
var mutex sync.Mutex

func HandleUdpConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	if conns == nil {
		plog.Println("Alloc Udp Conn Map Failed!")
		return
	}

	defer wg.Done()
	defer listener.udpconn.Close()
	defer func() {
		for _, v := range conns {
			select {
			case v <- nil:
			default:
			}
		}
	}()

	for {
		dgram := new(Dgram)
		if dgram == nil {
			plog.Println("Alloc Dgram Failed!")
			return
		}

		// Read Data From multi client.
		_, addr, err := DgramReadFromUdp(listener.udpconn, dgram, nil)
		if err != nil {
			plog.Println("Error : %s, Udp Listener %s Closed.", err.Error(),
				listener.udpconn.LocalAddr().String())
			return
		}

		// Dial To remote or always dialed.
		mutex.Lock()
		ch, ok := conns[addr.String()]
		mutex.Unlock()
		if !ok {
			// DialUDP will not blocking. because no communication with remote.
			co, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   net.ParseIP(cfg.RemoteAddr),
				Port: cfg.RemorePort,
			})
			if err != nil {
				plog.Println("Udp Dial Failed : %s.", err.Error())
				continue
			}

			ch = make(chan *Dgram, 32)
			if ch == nil {
				plog.Println("Alloc Udp Channel Failed!")
				return
			}

			mutex.Lock()
			conns[addr.String()] = ch
			mutex.Unlock()

			go HandleUdpData(listener.udpconn, addr, co, ch)
		}

		// using select because ch <- xx may blocking.
		select {
		case ch <- dgram:
		default:
		}
	}
}

func HandleUdpData(lconn *net.UDPConn, laddr *net.UDPAddr, rconn *net.UDPConn, ch chan *Dgram) {
	defer rconn.Close()

	plog.Println("New UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		laddr.String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	TransUdpData(lconn, laddr, rconn, ch)

	plog.Println("Destory UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		laddr.String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())
}

func TransUdpData(lconn *net.UDPConn, laddr *net.UDPAddr, rconn *net.UDPConn, ch chan *Dgram) {
	var (
		dgram   *Dgram
		err     error
		n       int
		bytes   int
		timeout = time.Now()
	)

	for {
		// if no data in five minutes, destory goroutine.
		if bytes != 0 {
			bytes = 0
			timeout = time.Now()
		} else {
			dur := time.Now().Sub(timeout)
			if dur > time.Minute*1 {
				goto ERR
			}
		}

		select {
		case dgram = <-ch:
			if dgram == nil {
				goto ERR
			}

			n, err = DgramWrite(rconn, dgram, nil)
			if err != nil {
				goto ERR
			}

			bytes += n

		default:
			t := time.Now().Add(time.Microsecond * time.Duration(1))
			n, err = DgramRead(rconn, dgram, &t)
			if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
				goto ERR
			}

			bytes += n

			n, err = DgramWriteToUdp(lconn, laddr, dgram, nil)
			if err != nil {
				goto ERR
			}

			bytes += n
		}
	}

ERR:
	mutex.Lock()
	delete(conns, laddr.String())
	mutex.Unlock()

	if err != nil {
		plog.Println("Error : %s On UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(), laddr.String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}
