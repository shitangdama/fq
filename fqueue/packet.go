package fqueue

import "time"

type Packet struct {
	// request   http.Request
	virfinish uint64 // 虚拟完成时间
	reqtime   uint64 //这个参数没有用到 这个是指需要 运行多少时间
	queue     *Queue
	starttime uint64
	endtime   uint64
	//
	key uint64
	seq uint64
	//
	estservicetime uint64 // 预估运行时间
	actservicetime uint64
}

// 更新队列 更新用的
func (p *Packet) updateTimeQueued() {
	p.queue.lock.Lock()
	defer p.queue.lock.Unlock()

	// 队列开始事件 定义开始事件

	if len(p.queue.Packets) == 0 && p.queue.packetExecuting == 0 {
		// p.queue.virStart is ignored
		// queues.lastvirfinish is in the virtualpast
		p.queue.virstart = uint64(time.Now().UnixNano())
	}

	// 如果 队列已经生成，就用最后的时间赋值 初始时间

	if len(p.queue.Packets) == 0 && p.queue.packetExecuting > 0 {
		p.queue.virstart = p.queue.lastvirfinish
	}

	// 这里计算不太对 这里是
	// 先算了 整体结束时间，为当前packet的结束时间
	// estservicetime 是预估时间

	p.virfinish = (uint64(len(p.queue.Packets)+1))*p.estservicetime + p.queue.virstart

	// 如果 队列中有包， 这结束时间等于 队列结束时间

	if len(p.queue.Packets) > 0 {
		// last virtual finish time of the queue is the virtual finish
		// time of the last request in the queue
		p.queue.lastvirfinish = p.virfinish // this pkt is last pkt
	}
}

// 出队列

func (p *Packet) updateTimeDequeued() {
	p.queue.lock.Lock()
	defer p.queue.lock.Unlock()
	// 预估时间 加到virstart上
	p.queue.virstart += p.estservicetime
}

// used when Packet's request is filled to store actual service time
func (p *Packet) updateTimeFinished() {
	p.queue.lock.Lock()
	defer p.queue.lock.Unlock()

	// 当任务结束的时候
	// 可执行剪一
	// ENDING PACKET SERVICING
	p.queue.packetExecuting--

	p.endtime = uint64(time.Now().UnixNano())
	p.actservicetime = p.starttime - p.endtime

	// 这里的逻辑是 开始时间 往后推迟 实际时间剪预估时间

	S := p.actservicetime
	G := p.estservicetime
	p.queue.virstart -= (G - S)
}
