package main

import (
	"log"
	"sync"
	"time"
)

var peerGroups = struct {
	m map[string]*peers
	sync.RWMutex
}{m: make(map[string]*peers)}

type peers struct {
	group   string
	seen    map[string]time.Time
	timeout time.Duration

	announce chan string
	req      chan chan []string
	ticker   <-chan time.Time
}

// get / make peers group
func getPeers(group string) *peers {
	peerGroups.RLock()
	if p, ok := peerGroups.m[group]; ok {
		peerGroups.RUnlock()
		return p
	}
	peerGroups.RUnlock()

	peerGroups.Lock()
	p := newPeers(group, purgeTime)
	peerGroups.m[group] = p
	peerGroups.Unlock()
	return p
}

func newPeers(group string, timeout time.Duration) *peers {
	p := &peers{
		group:   group,
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
	defer p.die()
	for {
		select {
		case addr := <-p.announce:
			log.Println("announce:", addr)
			p.seen[addr] = time.Now()
		case req := <-p.req:
			seen := p.list()
			req <- seen
		case now := <-p.ticker:
			p.purge(now)
			if len(p.seen) == 0 {
				// die
				return
			}
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

// cleanup
func (p *peers) die() {
	peerGroups.Lock()
	delete(peerGroups.m, p.group)
	peerGroups.Unlock()
}
