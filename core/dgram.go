package core

import (
	"net"
	"time"
)

const (
	DGRAM_SZ = 65536
)

type Dgram struct {
	Len    int
	buffer [DGRAM_SZ]byte
}

func DgramReadFromUdp(conn *net.UDPConn, dgram *Dgram, t *time.Time) (int, *net.UDPAddr, error) {
	if t != nil {
		conn.SetReadDeadline(*t)
	}

	n, addr, err := conn.ReadFromUDP(dgram.buffer[:])

	dgram.Len = n

	return n, addr, err
}

func DgramWriteToUdp(conn *net.UDPConn, addr *net.UDPAddr, dgram *Dgram, t *time.Time) (int, error) {
	if dgram.Len == 0 {
		return 0, nil
	}

	if t != nil {
		conn.SetWriteDeadline(*t)
	}

	n, err := conn.WriteToUDP(dgram.buffer[:dgram.Len], addr)

	return n, err
}

func DgramRead(conn *net.UDPConn, dgram *Dgram, t *time.Time) (int, error) {

	if t != nil {
		conn.SetReadDeadline(*t)
	}
	n, err := conn.Read(dgram.buffer[:])

	dgram.Len = n

	return n, err
}

func DgramWrite(conn *net.UDPConn, dgram *Dgram, t *time.Time) (int, error) {
	if dgram.Len == 0 {
		return 0, nil
	}

	if t != nil {
		conn.SetWriteDeadline(*t)
	}

	n, err := conn.Write(dgram.buffer[:dgram.Len])

	return n, err
}
