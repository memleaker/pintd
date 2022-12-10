package core

import (
	"net"
)

const (
	BUFFERSZ = 524288 // 2 ^ 19
)

type RingStream struct {
	In     uint32
	Out    uint32
	Size   uint32
	buffer [BUFFERSZ]byte
}

func NewRingStream() *RingStream {
	return &RingStream{
		In:   0,
		Out:  0,
		Size: BUFFERSZ,
	}
}

func (stream *RingStream) IsFull() bool {
	return ((stream.In - stream.Out) == stream.Size)
}

func (stream *RingStream) IsEmpty() bool {
	return (stream.In == stream.Out)
}

func min(a, b uint32) uint32 {
	if a > b {
		return b
	}

	return a
}

func (stream *RingStream) StreamRead(conn net.Conn) error {

	inpos := stream.In & (stream.Size - 1)
	free := stream.Size - (stream.In - stream.Out)
	rlen := stream.Size - inpos

	len := min(free, rlen)

	n, err := conn.Read(stream.buffer[inpos : inpos+len])

	stream.In += uint32(n)

	return err
}

func (stream *RingStream) StreamWrite(conn net.Conn) error {

	outpos := stream.Out & (stream.Size - 1)
	used := stream.In - stream.Out
	rlen := stream.Size - outpos

	len := min(used, rlen)

	n, err := conn.Write(stream.buffer[outpos : outpos+len])

	stream.Out += uint32(n)

	return err
}
