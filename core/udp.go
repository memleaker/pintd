package core

import (
	"errors"
	"net"
	"os"
	"pintd/config"
	"pintd/filter"
	"pintd/plog"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ConnInfo struct {
	Addr *net.UDPAddr
	Conn *net.UDPConn
}

// because UDP dont't have connection state. so we cannot know remote is or not closed.
// and it cause we cannot timely close expire goroutine when using one conn one goroutine model.
func HandleUdpConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	defer listener.udpconn.Close()

	var (
		conns sync.Map
		queue = NewDgramQueue()
		raddr = net.UDPAddr{IP: net.ParseIP(cfg.RemoteAddr), Port: cfg.RemorePort}
	)

	defer func() {
		conns.Range(func(key, value any) bool {
			rconn := value.(ConnInfo).Conn
			rconn.Close()
			return false
		})
	}()

	// read right multi server data to message queue.
	go ReadMultiServer(listener.udpconn, queue, &conns)

	// handle write data from message queue.
	go Write(listener.udpconn, queue)

	// read left multi client data to message queue.
	if err := ReadMultiClient(listener.udpconn, queue, &raddr, &conns, cfg); err != nil {
		plog.Println("Error : %s, Udp Listener %s Closed.", err.Error(),
			listener.udpconn.LocalAddr().String())
		return
	}
}

func ReadMultiClient(lconn *net.UDPConn, queue *DgramQueue, dial *net.UDPAddr, conns *sync.Map, cfg *config.RedirectConfig) error {
	var (
		rconn *net.UDPConn
		buf   = make([]byte, 65536)
	)

	for {
		n, laddr, err := lconn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// filter address
		ip, _, _ := strings.Cut(laddr.String(), ":")
		if deny := filter.DenyAccess(ip, cfg.SectionName); deny {
			plog.Println("Matched Deny Address : %s.", ip)
			continue
		}

		val, ok := conns.Load(laddr.String())
		if ok {
			rconn = val.(ConnInfo).Conn
		} else {
			rconn, err = net.DialUDP("udp", nil, dial)
			if err != nil {
				plog.Println("DialUDP Failed : %s", err.Error())
				continue
			}

			conninfo := ConnInfo{Addr: laddr, Conn: rconn}
			conns.Store(laddr.String(), conninfo)

			plog.Println("New UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
				laddr.String(), lconn.LocalAddr().String(),
				rconn.LocalAddr().String(), rconn.RemoteAddr().String())
		}

		PutDgramToQueue(queue, buf, n, rconn, dial)
	}
}

func ReadMultiServer(lconn *net.UDPConn, queue *DgramQueue, conns *sync.Map) {
	loop := true
	buf := make([]byte, 65536)

	for loop {
		conns.Range(func(key, value any) bool {
			rconn := value.(ConnInfo).Conn
			laddr := value.(ConnInfo).Addr

			rconn.SetReadDeadline(time.Now().Add(time.Microsecond * 1))
			n, err := rconn.Read(buf)
			if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
				conns.Delete(key)
				rconn.Close()
				plog.Println("Error : %s On UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
					err.Error(), laddr.String(), lconn.LocalAddr().String(),
					rconn.LocalAddr().String(), rconn.RemoteAddr().String())
				loop = false
				return false
			}

			if !errors.Is(err, os.ErrDeadlineExceeded) {
				PutDgramToQueue(queue, buf, n, lconn, laddr)
			}

			return true
		})

		runtime.Gosched()
	}
}

func Write(lconn *net.UDPConn, queue *DgramQueue) {
	var (
		err error
		buf = make([]byte, 65536)
	)

	for {
		// we don't need lock, because no race data.
		conn, addr, n := GetDgramFromQueue(queue, buf)

		// UDP don't have write buffer, write will not blocking.
		// if send failed, datadgram is lost.
		if conn == lconn {
			_, err = conn.WriteToUDP(buf[:n], addr) // send left
		} else {
			_, err = conn.Write(buf[:n]) // send right
		}

		if err != nil {
			plog.Println("Write Udp : %s", err.Error())
			return
		}
	}
}
