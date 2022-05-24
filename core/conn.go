package core

import (
	"errors"
	"io"
	"net"
	"os"
	"pintd/plog"
)

func HandleData(lconn, rconn net.Conn, ch chan int) {
	plog.Println("New Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.LocalAddr().String(), lconn.RemoteAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	defer lconn.Close()
	defer rconn.Close()

	TransData(lconn, rconn)

	plog.Println("Destory Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
		lconn.LocalAddr().String(), lconn.RemoteAddr().String(),
		rconn.LocalAddr().String(), rconn.RemoteAddr().String())

	ch <- CONN_DEC
}

func TransData(lconn, rconn net.Conn) {
	var (
		err     error
		lstream = new(Stream)
		rstream = new(Stream)
	)

	for {
		_, err = StreamRead(lconn, lstream)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}

		_, err = StreamWrite(rconn, lstream)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}

		_, err = StreamRead(rconn, rstream)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}

		_, err = StreamWrite(lconn, rstream)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			goto ERR
		}
	}

ERR:
	if err != io.EOF && err != io.ErrClosedPipe {
		plog.Println("Error : %s On Redirect Connection from [%s]->[%s] redirect to [%s]->[%s].",
			err.Error(),
			lconn.LocalAddr().String(), lconn.RemoteAddr().String(),
			rconn.LocalAddr().String(), rconn.RemoteAddr().String())
	}
}
