package core

import (
	"errors"
	"io"
	"net"
	"os"
	"pintd/config"
	"pintd/filter"
	"pintd/plog"
	"strconv"
	"strings"
	"sync"
	"time"
)

func HandleTcpConn(listener Listener, cfg *config.RedirectConfig, wg *sync.WaitGroup) {
	var (
		conns int = 0
		lconn net.Conn
		rconn net.Conn
		err   error
		ch    = make(chan int, 16)
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

		plog.Println("%s Accept Connection from %s.", lconn.LocalAddr().String(),
			lconn.RemoteAddr().String())

		// Dial to remote.
		for interval := 2; ; interval += 2 {
			if interval > 8 {
				plog.Println("Tcp Dial Failed : %s, Stop Reconnect.", err.Error())
				break
			}

			rconn, err = net.DialTimeout("tcp", cfg.RemoteAddr+":"+strconv.Itoa(cfg.RemorePort),
				time.Second*time.Duration(interval))
			if err != nil {
				plog.Println("Tcp Dial Failed : %s, Reconnect...", err.Error())
				continue
			}

			// Dial Success.
			conns++
			plog.Println("Tcp Dial to %s Success.", rconn.RemoteAddr().String())

			// handle data.
			go HandleTcpData(lconn, rconn, ch)
			break
		}
	}
}

func HandleTcpData(lconn, rconn net.Conn, ch chan int) {
	plog.Println("New TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	TransTcpData(lconn, rconn)

	plog.Println("Destory TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	// channel may blocking, so close connection first.
	lconn.Close()
	rconn.Close()

	ch <- CONN_DEC
}

func TransTcpData(lconn, rconn net.Conn) {
	var (
		err     error
		lstream = new(Stream)
		rstream = new(Stream)
	)

	if lstream == nil || rstream == nil {
		plog.Println("Alloc Stream Failed, Connection Closed.")
		goto ERR
	}

	for {
		_, err = StreamRead(lconn, lstream, time.Microsecond*time.Duration(500))
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			// Send residual data
			StreamWrite(rconn, lstream, time.Second*time.Duration(30))
			goto ERR
		}

		_, err = StreamWrite(rconn, lstream, time.Microsecond*time.Duration(500))
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}

		_, err = StreamRead(rconn, rstream, time.Microsecond*time.Duration(500))
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			// Send residual data
			StreamWrite(lconn, rstream, time.Second*time.Duration(30))
			goto ERR
		}

		_, err = StreamWrite(lconn, rstream, time.Microsecond*time.Duration(500))
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}
	}

ERR:
	if err != nil && err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Error : %s On TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(),
			lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}
