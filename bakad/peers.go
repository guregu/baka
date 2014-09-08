package main

import (
	"log"
	"time"
)

type peers struct {
	seen    map[string]time.Time
	timeout time.Duration

	announce chan string
	req      chan chan []string
	ticker   <-chan time.Time
}

func newPeers(timeout time.Duration) *peers {
	p := &peers{
		seen:    make(map[string]time.Time),
		timeout: timeout,

		announce: make(chan string),
		req:      make(chan chan []string),
		ticker:   time.Tick(5 * time.Second),
	}
	go p.run()
	return p
}

func (p *peers) run() {
	for {
		select {
		case addr := <-p.announce:
			log.Println("announce:", addr)
			p.seen[addr] = time.Now()
		case req := <-p.req:
			seen := p.list()
			go func() {
				req <- seen
			}()
		case now := <-p.ticker:
			p.purge(now)
		}
	}
}

func (p *peers) list() []string {
	seen := make([]string, 0, len(p.seen))
	for addr, _ := range p.seen {
		seen = append(seen, addr)
	}
	return seen
}

func (p *peers) purge(t time.Time) {
	for addr, last := range p.seen {
		if t.Sub(last) >= p.timeout {
			log.Println("purging dead peer", addr, "last seen", last.String())
			delete(p.seen, addr)
		}
	}
}
