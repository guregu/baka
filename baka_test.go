package baka

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/groupcache"
)

func TestBaka(t *testing.T) {
	rand.Seed(int64(time.Now().Nanosecond()))
	port := rand.Intn(9999) + 1000
	self := fmt.Sprintf("http://localhost:%d", port)
	server := "http://localhost:1337"
	pool := groupcache.NewHTTPPool(self)
	b := Baka{
		Server: server,
		Self:   self,
		Pool:   pool,
		Group:  "test",
		Rate:   time.Second * 10,
	}
	b.Run()
	time.Sleep(time.Minute)
}
