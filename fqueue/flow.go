package fqueue

import (
	"math/rand"
	"sync"
)

// 这里
type flowDesc struct {
	// In
	ftotal     uint64 // Total units in flow
	minreqtime uint64 // Min Packet reqtime
	maxreqtime uint64 // Max Packet reqtime

	// 这里少了一个权重
	// Out
	idealPercent  float64
	actualPercent float64
}

// 得到一个flow 流 这里要重写
func genFlow(fq *fqscheduler, desc *flowDesc, key uint64, done_wg *sync.WaitGroup) {
	for i, t := uint64(1), uint64(0); t < desc.ftotal; i++ {
		//time.Sleep(time.Microsecond)
		it := new(Packet)
		it.estservicetime = 100

		it.key = key

		// 初始化时间 最大最小 这个时间是什么意思

		//
		// 这里是一个设置 的属性，并不是一个赋值的属性，用于传进来的包的范围设置

		if desc.minreqtime == desc.maxreqtime {
			it.reqtime = desc.maxreqtime
		} else {
			it.reqtime = desc.minreqtime + uint64(rand.Int63n(int64(desc.maxreqtime-desc.minreqtime)))
		}

		if t+it.reqtime > desc.ftotal {
			it.reqtime = desc.ftotal - t
		}
		t += it.reqtime

		// seq表示包的序号
		it.seq = i
		// new packet
		fq.enqueue(it)
	}
	(*done_wg).Done()
}
