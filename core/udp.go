package core

import (
	"errors"
	"net"
	"os"
	"pintd/config"
	"pintd/plog"
	"runtime"
	"sync"
	"time"
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
		queue = NewDgramQueue()
		raddr = net.UDPAddr{IP: net.ParseIP(cfg.RemoteAddr), Port: cfg.RemorePort}
	)

	// read right multi server data to message queue.
	go ReadMultiServer(listener.udpconn, queue, &conns)

	// handle write data from message queue.
	go Write(listener.udpconn, queue)

	// read left multi client data to message queue.
	for {
		err := ReadMultiClient(listener.udpconn, queue, &raddr, &conns)
		if err != nil {
			plog.Println("Error : %s, Udp Listener %s Closed.", err.Error(),
				listener.udpconn.LocalAddr().String())
			return
		}
	}
}

func ReadMultiClient(lconn *net.UDPConn, queue *DgramQueue, dial *net.UDPAddr, conns *sync.Map) error {
	var (
		rconn *net.UDPConn
		buf   = make([]byte, 65536)
	)

	n, laddr, err := lconn.ReadFromUDP(buf)
	if err != nil {
		return err
	}

	val, ok := conns.Load(laddr.String())
	if ok {
		rconn = val.(ConnInfo).Conn
	} else {
		rconn, err = net.DialUDP("udp", nil, dial)
		if err != nil {
			plog.Println("DialUDP Failed : %s", err.Error())
			return nil
		}

		conninfo := ConnInfo{Addr: laddr, Conn: rconn}
		conns.Store(laddr.String(), conninfo)

		plog.Println("New UDP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			laddr.String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}

	PushDgramToQueue(queue, buf, n, rconn, dial)

	return nil
}

func ReadMultiServer(lconn *net.UDPConn, queue *DgramQueue, conns *sync.Map) {

	buf := make([]byte, 65536)

	for {
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
				return false
			}

			if !errors.Is(err, os.ErrDeadlineExceeded) {
				PushDgramToQueue(queue, buf, n, lconn, laddr)
			}

			return true
		})

		runtime.Gosched()
	}
}

func Write(lconn *net.UDPConn, queue *DgramQueue) {

	for {
		// we don't need lock, because no race data.
		conn, addr, data := GetDgramFromQueue(queue)

		// UDP don't have write buffer, write will not blocking.
		// if send failed, datadgram is lost.
		if conn == lconn {
			conn.WriteToUDP(data, addr) // send left
		} else {
			conn.Write(data) // send right
		}
	}
}
