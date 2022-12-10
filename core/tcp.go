package core

import (
	"errors"
	"fmt"
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

	defer wg.Done()
	defer listener.listener.Close()
	defer func() {
		if err := recover(); err != nil {
			plog.Println(fmt.Sprint(err))
		}
	}()

	conns := make(chan bool, cfg.MaxRedirects)

	// Wait tcp connection coming.
	for {
		lconn, err := listener.listener.Accept()
		if err != nil {
			plog.Println("Accept Connection Failed %s, listener closed.", err.Error())
			return
		}

		// Limit ip address
		ip, _, _ := strings.Cut(lconn.RemoteAddr().String(), ":")
		if filter.DenyAccess(ip, cfg.SectionName) {
			lconn.Close()
			plog.Println("Matched Deny Address : %s.", ip)
			continue
		}

		// Limit the number of connections
		select {
		case conns <- true:
		default:
			lconn.Close()
			plog.Println("Connection Limit to %d, Closed Connection.", cfg.MaxRedirects)
			continue
		}

		plog.Println("%s Accept Connection from %s.", lconn.LocalAddr().String(),
			lconn.RemoteAddr().String())

		// Dial to remote.
		go DialToRemote(lconn, cfg, conns)
	}
}

func DialToRemote(lconn net.Conn, cfg *config.RedirectConfig, conns chan bool) {

	// decrease conn number
	defer func() { <-conns }()
	defer func() {
		if err := recover(); err != nil {
			plog.Println(fmt.Sprint(err))
		}
	}()

	// dial to remote
	rconn, err := net.DialTimeout("tcp", cfg.RemoteAddr+":"+strconv.Itoa(cfg.RemorePort),
		time.Second*time.Duration(3))
	if err != nil {
		lconn.Close()
		plog.Println("Tcp Dial Failed : %s", err.Error())
		return
	}

	plog.Println("Tcp Dial to %s Success.", rconn.RemoteAddr().String())

	// set tcp keepalive 60 seconds
	SetTcpKeepalive(lconn, 60)
	SetTcpKeepalive(rconn, 60)

	// handle data.
	HandleTcpData(lconn, rconn, conns, cfg)
}

func HandleTcpData(lconn, rconn net.Conn, conns chan bool, cfg *config.RedirectConfig) {
	var wg sync.WaitGroup

	plog.Println("New TCP Redirect Connection [%s]<->[%s] redirect [%s]<->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	SetTcpNodelay(lconn, cfg.NoDelay)
	SetTcpNodelay(rconn, cfg.NoDelay)

	if cfg.NoDelay {
		RedirectNodelay(lconn, rconn, &wg)
	} else {
		Redirect(lconn, rconn, &wg)
	}

	wg.Wait()

	plog.Println("Destory TCP Redirect Connection [%s]<->[%s] redirect [%s]<->[%s].",
		lconn.RemoteAddr().String(), lconn.LocalAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())
}

func RedirectNodelay(lconn, rconn net.Conn, wg *sync.WaitGroup) {

	go RedirectIo(lconn, rconn, wg)
	go RedirectIo(rconn, lconn, wg)

	wg.Add(2)
}

func Redirect(lconn, rconn net.Conn, wg *sync.WaitGroup) {
	lstream := NewRingStream()
	rstream := NewRingStream()

	go LeftRead(lconn, rconn, lstream, wg)
	go RightWrite(lconn, rconn, lstream, wg)

	go RightRead(lconn, rconn, rstream, wg)
	go LeftWrite(lconn, rconn, rstream, wg)

	wg.Add(4)
}

func LeftRead(lconn, rconn net.Conn, lstream *RingStream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	RedirectRead(lconn, lstream)
}

func RightRead(lconn, rconn net.Conn, rstream *RingStream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	RedirectRead(rconn, rstream)
}

func RedirectRead(conn net.Conn, stream *RingStream) {
	var err error

	for {
		err = stream.StreamRead(conn)
		if err != nil {
			break
		}
	}

	if !errors.Is(err, net.ErrClosed) &&
		err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Tcp Redirect Error : %s", err.Error())
	}
}

func LeftWrite(lconn, rconn net.Conn, rstream *RingStream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	RedirectWrite(lconn, rstream)
}

func RightWrite(lconn, rconn net.Conn, lstream *RingStream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	RedirectWrite(rconn, lstream)
}

func RedirectWrite(conn net.Conn, stream *RingStream) {

	var err error

	for {
		err = stream.StreamWrite(conn)
		if err != nil {
			break
		}
	}

	if !errors.Is(err, net.ErrClosed) &&
		err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Tcp Redirect Error : %s", err.Error())
	}
}

func RedirectIo(lconn, rconn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer lconn.Close()
	defer rconn.Close()

	var (
		err error
		n   int
		buf = make([]byte, BUFFERSZ)
	)

	for {
		n, err = lconn.Read(buf)
		if err != nil {
			// Send residual data
			rconn.Write(buf[:n])
			break
		}

		// conn.Write will block wait until all data copid to kernel
		// so we dont't need stream or ringbuffer to cache data
		// unless set deadline, conn.SetDeadline()
		_, err = rconn.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// net.ErrClosed : IO call on a closed network conn or closed by another
	//   goroutine before IO is complate.
	// io.EOF : Read EOF, no more input is available
	// io.ErrClosedPipe : read or write on a closed pipe
	if !errors.Is(err, net.ErrClosed) &&
		err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Tcp Redirect Error : %s", err.Error())
	}
}

func SetTcpKeepalive(conn net.Conn, secs int) error {
	tconn, ok := conn.(*net.TCPConn)
	if !ok {
		return errors.New("invalid argument")
	}

	if err := tconn.SetKeepAlive(true); err != nil {
		return err
	}

	if err := tconn.SetKeepAlivePeriod(time.Duration(secs) * time.Second); err != nil {
		return err
	}

	return nil
}

func SetTcpNodelay(conn net.Conn, delay bool) error {
	tconn, ok := conn.(*net.TCPConn)
	if !ok {
		return errors.New("invalid argument")
	}

	return tconn.SetNoDelay(delay)
}
