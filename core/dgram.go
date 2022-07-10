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
	Sem        chan bool
	GetIndex   int
	PutIndex   int
	PutMutex   sync.Mutex
	GetMutex   sync.Mutex
	DgramSlice [QUEUE_SZ]Dgram
	Buffer     [QUEUE_SZ][DGRAM_SZ]byte
}

func NewDgramQueue() *DgramQueue {
	var queue DgramQueue

	queue.Sem = make(chan bool, QUEUE_SZ)

	return &queue
}

func PushDgramToQueue(queue *DgramQueue, buf []byte, n int, conn *net.UDPConn, dstaddr *net.UDPAddr) {

	queue.Sem <- true

	queue.PutMutex.Lock()

	index := queue.PutIndex
	queue.PutIndex = (queue.PutIndex + 1) % QUEUE_SZ

	queue.PutMutex.Unlock()

	copy(queue.Buffer[index][:], buf[:n])

	queue.DgramSlice[index].Addr = dstaddr
	queue.DgramSlice[index].Conn = conn
	queue.DgramSlice[index].Data = queue.Buffer[index][:n]
}

func GetDgramFromQueue(queue *DgramQueue) (*net.UDPConn, *net.UDPAddr, []byte) {

	select {
	case <-queue.Sem:
		queue.GetMutex.Lock()

		index := queue.GetIndex
		queue.GetIndex = (queue.GetIndex + 1) % QUEUE_SZ

		queue.GetMutex.Unlock()

		return queue.DgramSlice[index].Conn, queue.DgramSlice[index].Addr, queue.DgramSlice[index].Data
	}
}
