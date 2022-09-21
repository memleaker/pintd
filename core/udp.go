package core

import (
	"errors"
	"net"
	"pintd/config"
	"pintd/filter"
	"pintd/plog"
	"strings"
	"sync"
)

type ConnInfo struct {
	Addr *net.UDPAddr
	Conn *net.UDPConn
}

func HandleUdpConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	defer listener.udpconn.Close()

	var (
		conns sync.Map
		rconn *net.UDPConn
		buf   = make([]byte, 65535)
		raddr = net.UDPAddr{IP: net.ParseIP(cfg.RemoteAddr), Port: cfg.RemorePort}
	)

	// close all conn.
	defer func() {
		conns.Range(func(key, value any) bool {
			rconn := value.(ConnInfo).Conn
			rconn.Close()
			return false
		})
	}()

	// read left multi client data.
	for {
		n, laddr, err := listener.udpconn.ReadFromUDP(buf)
		if err != nil {
			plog.Println("Error : %s, Udp Listener %s Closed.", err.Error(),
				listener.udpconn.LocalAddr().String())
			return
		}

		// filter address
		ip, _, _ := strings.Cut(laddr.String(), ":")
		if deny := filter.DenyAccess(ip, cfg.SectionName); deny {
			plog.Println("Matched Deny Address : %s.", ip)
			continue
		}

		// new connection.
		val, ok := conns.Load(laddr.String())
		if ok {
			rconn = val.(ConnInfo).Conn
		} else {
			rconn, err = net.DialUDP("udp", nil, &raddr)
			if err != nil {
				plog.Println("DialUDP Failed : %s", err.Error())
				continue
			}

			conninfo := ConnInfo{Addr: laddr, Conn: rconn}
			conns.Store(laddr.String(), conninfo)

			go UdpRightToLeft(listener.udpconn, rconn, laddr, &conns)

			plog.Println("New UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
				laddr.String(), listener.udpconn.LocalAddr().String(),
				rconn.LocalAddr().String(), rconn.RemoteAddr().String())
		}

		// UDP don't have write buffer, write will not blocking.
		// if send failed, datadgram is lost.
		rconn.Write(buf[:n])
	}
}

func UdpRightToLeft(lconn, rconn *net.UDPConn, laddr *net.UDPAddr, conns *sync.Map) {
	buf := make([]byte, 65536)

	defer conns.Delete(laddr.String())
	defer rconn.Close()

	for {
		n, err := rconn.Read(buf)
		if err != nil {
			// conn closed
			if errors.Is(err, net.ErrClosed) {
				return
			}

			plog.Println("Error : %s On UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
				err.Error(), laddr.String(), lconn.LocalAddr().String(),
				rconn.LocalAddr().String(), rconn.RemoteAddr().String())
			return
		}

		lconn.WriteToUDP(buf[:n], laddr)
	}
}
