package main

import (
	"testing"
	"time"
)

func TestPurge(t *testing.T) {
	p := newTestPeers()
	p.purge(time.Now())
	// http://dead.com should get deleted
	if _, exists := p.seen["http://dead.com"]; exists {
		t.Error("dead peer was not purged")
	}
	if _, exists := p.seen["http://alive.com"]; !exists {
		t.Error("alive peer was purged")
	}
}

func newTestPeers() *peers {
	return &peers{
		seen: map[string]time.Time{
			"http://alive.com": time.Now(),
			"http://dead.com":  time.Now().Add(-2 * time.Minute),
		},
		timeout: time.Minute,
	}
}
