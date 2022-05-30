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

func HandleUdpConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	if conns == nil {
		plog.Println("Alloc Udp Conn Map Failed!")
		return
	}

	defer wg.Done()
	defer listener.udpconn.Close()
	defer func() {
		for _, v := range conns {
			v <- nil
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
		ch, ok := conns[addr.String()]
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

			ch = make(chan *Dgram, 8)
			if ch == nil {
				plog.Println("Alloc Udp Channel Failed!")
				return
			}

			conns[addr.String()] = ch

			go HandleUdpData(listener.udpconn, addr, co, ch)
		}

		ch <- dgram
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
		dgram *Dgram
		err   error
	)

	for {
		select {
		case dgram = <-ch:
			if dgram == nil {
				goto ERR
			}

			_, err = DgramWrite(rconn, dgram, nil)
			if err != nil {
				goto ERR
			}

		default:
			t := time.Now().Add(time.Millisecond * time.Duration(3))
			_, err = DgramRead(rconn, dgram, &t)
			if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
				goto ERR
			}

			_, err = DgramWriteToUdp(lconn, laddr, dgram, nil)
			if err != nil {
				goto ERR
			}
		}
	}

ERR:
	delete(conns, laddr.String())
	if err != nil {
		plog.Println("Error : %s On UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(), laddr.String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}
