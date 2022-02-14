package fqueue

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func consumeQueue(t *testing.T, fq *fqscheduler, descs []flowDesc) (float64, error) {
	// 活跃queue
	active := make(map[uint64]bool)
	var total uint64
	acnt := make(map[uint64]uint64)
	cnt := make(map[uint64]uint64)
	seqs := make(map[uint64]uint64)

	// wsum是什么参数
	// 是不是weight总数量

	var wsum uint64
	for range descs {
		wsum += uint64(0 + 1)
	}

	// 这个参数没用到
	// pktTotal := uint64(0)
	// for _, desc := range descs {
	// 	pktTotal += desc.ftotal
	// }

	// 不停的把 把逻辑取出来

	// stdDev appears to change quite a bit if the queue is dequeued without
	// waiting
	time.Sleep(1 * time.Second)
	for i, ok := fq.dequeue(); ok; i, ok = fq.dequeue() {

		// 真实需要这样运行
		time.Sleep(time.Duration(i.reqtime) * time.Nanosecond) // Simulate constrained bandwidth
		// time.Sleep(time.Microsecond) // Simulate constrained bandwidth
		// TODO(aaron-prindle) added this
		i.updateTimeFinished()

		// 这里判断之前的包 + 1 = 后面包的seq
		it := i
		seq := seqs[it.key]
		if seq+1 != it.seq {
			return 0, fmt.Errorf("Packet for flow %d came out of queue out-of-order: expected %d, got %d", it.key, seq+1, it.seq)
		}
		seqs[it.key] = it.seq

		// cnt是什么意思， 是总运行时间

		if cnt[it.key] == 0 {
			active[it.key] = true
		}
		cnt[it.key] += it.reqtime

		if len(active) == len(descs) {
			acnt[it.key] += it.reqtime
			total += it.reqtime
		}

		if cnt[it.key] == descs[it.key].ftotal {
			delete(active, it.key)
		}

	}

	var variance float64
	for key := uint64(0); key < uint64(len(descs)); key++ {
		// if total is 0, percents become NaN and those values would otherwise pass
		if total == 0 {
			t.Fatalf("expected 'total' to be nonzero")
		}
		descs[key].idealPercent = (((float64(total) * float64(0+1)) / float64(wsum)) / float64(total)) * 100
		descs[key].actualPercent = (float64(acnt[key]) / float64(total)) * 100
		x := descs[key].idealPercent - descs[key].actualPercent
		x *= x
		variance += x
	}

	stdDev := math.Sqrt(variance)
	return stdDev, nil
}
