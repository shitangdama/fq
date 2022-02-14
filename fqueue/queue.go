package fqueue

import "sync"

type Queue struct {
	lock          sync.Mutex
	Packets       []*Packet
	key           uint64
	lastvirfinish uint64 // 最后结束时间
	//
	packetExecuting uint64 // 执行数量
	virstart        uint64 // 队列开始时间
}

func (q *Queue) enqueue(packet *Packet) {
	q.Packets = append(q.Packets, packet)
}

func (q *Queue) dequeue() (*Packet, bool) {
	if len(q.Packets) == 0 {
		return nil, false
	}
	packet := q.Packets[0]
	q.Packets = q.Packets[1:]
	q.packetExecuting++
	return packet, true
}

// 初始化队列
func InitQueues(n int, key uint64) []*Queue {
	queues := []*Queue{}
	for i := 0; i < n; i++ {
		qkey := key
		if key == 0 {
			qkey = uint64(i)
		}
		queues = append(queues, &Queue{
			Packets: []*Packet{},
			key:     qkey,
		})
	}

	return queues
}
