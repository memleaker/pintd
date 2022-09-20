package core

import (
	"errors"
	"io"
	"net"
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
		lconn net.Conn
		rconn net.Conn
		err   error
		ch    = make(chan bool, cfg.MaxRedirects)
	)

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

		plog.Println("%s Accept Connection from %s.", lconn.LocalAddr().String(),
			lconn.RemoteAddr().String())

		// Dial to remote.
		rconn, err = net.DialTimeout("tcp", cfg.RemoteAddr+":"+strconv.Itoa(cfg.RemorePort),
			time.Second*time.Duration(3))
		if err != nil {
			lconn.Close()
			plog.Println("Tcp Dial Failed : %s", err.Error())
			continue
		}

		// Dial Success.
		// check connections number.
		select {
		case ch <- true:
		default:
			lconn.Close()
			plog.Println("Connection Limit to %d, Closed Connection.", cfg.MaxRedirects)
			continue
		}
		plog.Println("Tcp Dial to %s Success.", rconn.RemoteAddr().String())

		// handle data.
		go HandleTcpData(lconn, rconn, ch)
	}
}

func HandleTcpData(lconn, rconn net.Conn, ch chan bool) {
	var (
		wg      sync.WaitGroup
		lstream = NewStream()
		rstream = NewStream()
	)

	plog.Println("New TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	go LeftToRight(lconn, rconn, lstream, rstream, &wg)
	go RightToLeft(lconn, rconn, lstream, rstream, &wg)

	wg.Add(2)
	wg.Wait()

	plog.Println("Destory TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	// decrease channel.
	<-ch
}

func LeftToRight(lconn, rconn net.Conn, lstream, rstream *Stream, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	for {
		if err = StreamRead(lconn, lstream); err != nil {
			// Send residual data
			StreamWrite(rconn, lstream)
			goto ERR
		}

		if err = StreamWrite(rconn, lstream); err != nil {
			goto ERR
		}
	}

ERR:
	if err != nil && !errors.Is(err, net.ErrClosed) &&
		err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Error : %s On TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(),
			lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}

func RightToLeft(lconn, rconn net.Conn, lstream, rstream *Stream, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	for {
		if err = StreamRead(rconn, rstream); err != nil {
			// Send residual data
			StreamWrite(lconn, rstream)
			goto ERR
		}

		if err = StreamWrite(lconn, rstream); err != nil {
			goto ERR
		}
	}

ERR:
	if err != nil && !errors.Is(err, net.ErrClosed) &&
		err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Error : %s On TCP Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(),
			lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}
