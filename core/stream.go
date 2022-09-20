package core

import (
	"net"
)

const (
	BUFFERSZ = 655360
)

type Stream struct {
	Head   int
	Tail   int
	buffer [BUFFERSZ]byte
}

func NewStream() *Stream {
	return &Stream{
		Head: 0,
		Tail: 0,
	}
}

func StreamRead(conn net.Conn, stream *Stream) error {
	if stream.Tail == BUFFERSZ {
		return nil
	}

	n, err := conn.Read(stream.buffer[stream.Tail:])
	if err != nil {
		stream.Tail += n
		return err
	}

	stream.Tail += n

	return nil
}

func StreamWrite(conn net.Conn, stream *Stream) error {
	tail := stream.Tail

	for stream.Head != tail {
		n, err := conn.Write(stream.buffer[stream.Head:tail])
		if err != nil {
			// even write timeout, n may bigger than 0.
			stream.Head += n
			return err
		}

		stream.Head += n
	}

	// write all data ok, reset buffer.
	if tail == BUFFERSZ {
		stream.Head = 0
		stream.Tail = 0
	}

	return nil
}
