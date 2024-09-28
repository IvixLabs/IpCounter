package main

import (
	"math/bits"
	"sync"
)

type IP3Collector struct {
	ipChan chan IP3
	ips    [256 * 256 * 4]uint64
}

func NewIP3Collector(wg *sync.WaitGroup) *IP3Collector {
	iPCollector := &IP3Collector{
		ipChan: make(chan IP3),
	}

	wg.Add(1)
	go func(collector *IP3Collector) {
		defer wg.Done()

		for ip := range collector.ipChan {
			bucketIdx := int(ip[0])<<10 | int(ip[1])<<2 | int(ip[2])>>6

			mask := uint64(1)
			move := ip[2] & 0b111111
			mask = mask << move

			collector.ips[bucketIdx] = collector.ips[bucketIdx] | mask
		}

	}(iPCollector)

	return iPCollector
}

func (iPCollector *IP3Collector) Close() {
	close(iPCollector.ipChan)
}

func (iPCollector *IP3Collector) SendIP(ip IP3) {
	iPCollector.ipChan <- ip
}

func (iPCollector *IP3Collector) Len() int {
	cnt := 0
	for _, v := range iPCollector.ips {
		cnt += bits.OnesCount64(v)
	}
	return cnt
}

type IPCollectorProcessor struct {
	ipCollectors [256]*IP3Collector
	wg           *sync.WaitGroup
}

func NewIpCollectorProcessor() *IPCollectorProcessor {
	collectors := [256]*IP3Collector{}
	wg := &sync.WaitGroup{}

	for i := range 256 {
		collectors[i] = NewIP3Collector(wg)
	}

	return &IPCollectorProcessor{
		wg:           wg,
		ipCollectors: collectors,
	}
}

func (p *IPCollectorProcessor) HandleIP(ip IP) {
	ipCollIdx := ip[0]
	ipColl := p.ipCollectors[ipCollIdx]
	ipColl.SendIP(ip.ToIP3())
}

func (p *IPCollectorProcessor) Close() {
	for _, ipColl := range p.ipCollectors {
		ipColl.Close()
	}
}

func (p *IPCollectorProcessor) Wait() {
	p.wg.Wait()
}

func (p *IPCollectorProcessor) CalcUniqIPs() int {
	uniqIPs := 0
	for _, ipColl := range p.ipCollectors {
		uniqIPs += ipColl.Len()
	}

	return uniqIPs
}
