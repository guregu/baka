package baka

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/groupcache"
)

type Baka struct {
	Server string
	Self   string
	Group  string
	Pool   *groupcache.HTTPPool
	Rate   time.Duration

	ticker <-chan time.Time
}

func (b *Baka) Run() {
	if b.Server == "" || b.Self == "" ||
		b.Pool == nil || b.Rate == 0 {
		panic("baka: missing params")
	}

	b.ticker = time.Tick(b.Rate)
	b.update()
	go b.run()
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
	api := fmt.Sprintf("%s/%s/announce", b.Server, b.Group)
	log.Println("URL", api)
	resp, err := http.PostForm(api, url.Values{"url": {b.Self}})
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
	b.Pool.Set(peers...)
}
