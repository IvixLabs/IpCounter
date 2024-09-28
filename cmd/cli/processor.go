package main

import (
	"math/bits"
	"sync"
)

type Processor struct {
	pool    []chan IP
	poolLen int
	ips     []uint64
	wg      *sync.WaitGroup
}

func NewProcessor(poolLen int) *Processor {

	pool := make([]chan IP, poolLen)
	for i := range poolLen {
		pool[i] = make(chan IP, 1000)
	}

	ips := make([]uint64, 256*256*256*4)

	proc := &Processor{
		pool:    pool,
		poolLen: poolLen,
		ips:     ips,
		wg:      &sync.WaitGroup{},
	}
	proc.startLoop()

	return proc

}

func (p *Processor) HandleIP(ip IP) {
	poolIdx := int(ip[0]) / 256 * p.poolLen
	p.pool[poolIdx] <- ip

}

func (p *Processor) startLoop() {
	for _, ch := range p.pool {
		p.wg.Add(1)
		go func(ch chan IP) {
			for ip := range ch {
				p.updateBucket(ip)
			}
			p.wg.Done()
		}(ch)
	}
}

func (p *Processor) Close() {
	for _, ch := range p.pool {
		close(ch)
	}
}

func (p *Processor) Wait() {
	p.wg.Wait()
}

func (p *Processor) updateBucket(ip IP) {
	idx := int(ip[0])<<18 | int(ip[1])<<10 | int(ip[2])<<2 | int(ip[3])>>6

	mask := uint64(1)
	move := ip[3] & 0b111111
	mask = mask << move

	p.ips[idx] = p.ips[idx] | mask

}

func (p *Processor) CalcUniqIPs() int {
	cnt := 0
	for _, v := range p.ips {
		cnt += bits.OnesCount64(v)
	}

	return cnt
}
