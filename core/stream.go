package core

import (
	"net"
	"pintd/plog"
	"time"
)

const (
	BUFFERSZ = 655360
)

type Stream struct {
	Head   int
	Tail   int
	buffer [BUFFERSZ]byte
}

func StreamRead(conn net.Conn, stream *Stream, dur time.Duration) (int, error) {
	var num int = 0

	for stream.Tail != BUFFERSZ {
		if err := conn.SetReadDeadline(time.Now().Add(dur)); err != nil {
			plog.Println("Set ReadDeadline Failed %s.", err.Error())
			return 0, nil
		}

		n, err := conn.Read(stream.buffer[stream.Tail:])
		if err != nil {
			num += n
			stream.Tail += n
			return num, err
		}

		num += n
		stream.Tail += n

	}

	return num, nil
}

func StreamWrite(conn net.Conn, stream *Stream, dur time.Duration) (int, error) {
	var num int = 0

	if stream.Tail == 0 {
		return 0, nil
	}

	for stream.Head != stream.Tail {
		if err := conn.SetWriteDeadline(time.Now().Add(dur)); err != nil {
			plog.Println("Set WriteDeadline Failed %s.", err.Error())
			return 0, nil
		}

		n, err := conn.Write(stream.buffer[stream.Head:stream.Tail])
		if err != nil {
			// even write timeout, n may bigger than 0.
			num += n
			stream.Head += n
			return num, err
		}

		num += n
		stream.Head += n
	}

	stream.Head = 0
	stream.Tail = 0

	return num, nil
}
