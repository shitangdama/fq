package fqueue

import (
	"math"
	"sync"
	"time"
)

// 调度器
// 这里后面要加一下动态添加
// queue

type fqscheduler struct {
	lock   sync.Mutex
	queues []*Queue
	closed bool
}

func newfqscheduler(queues []*Queue) *fqscheduler {
	fq := &fqscheduler{
		queues: queues,
	}
	return fq
}

func (q *fqscheduler) chooseQueue(packet *Packet) *Queue {
	for _, queue := range q.queues {
		if packet.key == queue.key {
			// use shuffle sharding to get a queue
			packet.queue = queue
			return queue
		}
	}
	// fmt.Printf("packet w/ key: %v\n", packet.key)
	panic("no matching queue for packet")
}

func (q *fqscheduler) enqueue(packet *Packet) {
	q.lock.Lock()
	defer q.lock.Unlock()
	queue := q.chooseQueue(packet)
	queue.enqueue(packet)

	// STARTING PACKET SERVICING
	packet.starttime = uint64(time.Now().UnixNano())

	packet.updateTimeQueued()
}

// 这里出队列的逻辑要修改一下
func (q *fqscheduler) dequeue() (*Packet, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	queue := q.selectQueue()
	if queue == nil {
		return nil, false
	}
	packet, ok := queue.dequeue()
	if ok {
		packet.endtime = uint64(time.Now().UnixNano())
	}
	if ok {
		packet.updateTimeDequeued()
	}
	return packet, ok
}

func (q *fqscheduler) selectQueue() *Queue {
	minvirfinish := uint64(math.MaxUint64)
	var minqueue *Queue
	for _, queue := range q.queues {
		if len(queue.Packets) != 0 && queue.Packets[0].virfinish < minvirfinish {
			minvirfinish = queue.Packets[0].virfinish
			minqueue = queue
		}
	}
	return minqueue
}
