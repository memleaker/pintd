package core

import (
	"net"
	"sync"
)

const (
	DGRAM_SZ = 65536
	QUEUE_SZ = 64
)

type Dgram struct {
	Conn *net.UDPConn
	Addr *net.UDPAddr
	Data []byte
}

type DgramQueue struct {
	GetSem     chan bool
	PutSem     chan bool
	GetIndex   int
	PutIndex   int
	GetMutex   sync.Mutex
	PutMutex   sync.Mutex
	DgramSlice [QUEUE_SZ]Dgram
	Buffer     [QUEUE_SZ][DGRAM_SZ]byte
}

func NewDgramQueue() *DgramQueue {
	queue := new(DgramQueue)

	queue.GetSem = make(chan bool, QUEUE_SZ)
	queue.PutSem = make(chan bool, QUEUE_SZ)

	return queue
}

func PutDgramToQueue(queue *DgramQueue, buf []byte, n int, conn *net.UDPConn, dstaddr *net.UDPAddr) {

	// if queue is not empty, write and increase put sem.
	queue.PutSem <- true

	// lock and increase putindex.
	queue.PutMutex.Lock()
	index := queue.PutIndex
	queue.PutIndex = (queue.PutIndex + 1) % QUEUE_SZ
	queue.PutMutex.Unlock()

	// assignment.
	copy(queue.Buffer[index][:], buf[:n])
	queue.DgramSlice[index].Addr = dstaddr
	queue.DgramSlice[index].Conn = conn
	queue.DgramSlice[index].Data = queue.Buffer[index][:n]

	// push ok, increase Get Sem.
	queue.GetSem <- true
}

func GetDgramFromQueue(queue *DgramQueue, buf []byte) (*net.UDPConn, *net.UDPAddr, int) {

	// Sem is readable, message queue is not null.
	// read and decrease get sem.
	<-queue.GetSem

	// increase Getindex.
	queue.GetMutex.Lock()
	index := queue.GetIndex
	queue.GetIndex = (queue.GetIndex + 1) % QUEUE_SZ
	queue.GetMutex.Unlock()

	// assignment.
	len := copy(buf, queue.DgramSlice[index].Data)
	conn := queue.DgramSlice[index].Conn
	addr := queue.DgramSlice[index].Addr

	// get ok, decrease put sem.
	<-queue.PutSem

	return conn, addr, len
}
