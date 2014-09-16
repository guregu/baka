package baka

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/groupcache"
)

type Baka struct {
	server string
	self   string
	pool   *groupcache.HTTPPool

	ticker <-chan time.Time
}

func (b *Baka) run() {
	for {
		select {
		case <-b.ticker:
			b.update()
		}
	}
}

func (b *Baka) update() {
	resp, err := http.PostForm(b.server+"/announce",
		url.Values{"url": {b.self}})
	if err != nil {
		log.Println("baka error", err)
		return
	}
	var peers []string
	err = json.NewDecoder(resp.Body).Decode(&peers)
	if err != nil {
		log.Println("baka error", err)
		return
	}
	// log.Println("got peers", peers)
	b.pool.Set(peers...)
}

// Update starts listening for updates to the peer list. Server and self should be URLs.
func Update(server, self string, pool *groupcache.HTTPPool, announceRate time.Duration) {
	b := &Baka{
		server: server,
		self:   self,
		pool:   pool,

		ticker: time.Tick(announceRate),
	}
	b.update()
	go b.run()
}
