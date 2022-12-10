package core

import (
	"fmt"
	"net"
)

const (
	// BUFFERSZ = 524288 // 2 ^ 19
	BUFFERSZ = 16
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

func (stream *RingStream) StreamRead(conn net.Conn) error {
	if stream.IsFull() {
		return nil
	}

	inpos := stream.In & (stream.Size - 1)
	outpos := stream.Out & (stream.Size - 1)
	end := stream.Size

	if inpos < outpos {
		end = outpos
	}

	n, err := conn.Read(stream.buffer[inpos:end])

	stream.In += uint32(n)

	fmt.Println("recv ", n, "In", stream.In, "Out", stream.Out)

	return err
}

func (stream *RingStream) StreamWrite(conn net.Conn) error {
	if stream.IsEmpty() {
		return nil
	}

	inpos := stream.In & (stream.Size - 1)
	outpos := stream.Out & (stream.Size - 1)
	end := stream.Size

	if inpos >= outpos {
		end = inpos
	}

	n, err := conn.Write(stream.buffer[outpos:end])

	stream.Out += uint32(n)

	fmt.Println("w", n, "Out", stream.Out)

	return err
}
