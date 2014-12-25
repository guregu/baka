package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var bind = flag.String("bind", ":1337", "bind address")
var timeout = flag.String("timeout", "1m", "dead peer timeout")
var purgeTime time.Duration

func main() {
	flag.Parse()

	var err error
	purgeTime, err = time.ParseDuration(*timeout)
	if err != nil {
		panic(err)
	}

	setup()
}

func setup() {
	router := httprouter.New()
	router.GET("/:group/peers", peersHandler)
	router.POST("/:group/announce", announceHandler)

	log.Fatal(http.ListenAndServe(*bind, router))
}

func announceHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	group := params.ByName("group")
	if group == "" {
		http.Error(w, "no group", 400)
		return
	}
	peerlist := getPeers(group)

	if r.Method != "POST" {
		http.Error(w, "bad method", 400)
		return
	}
	url := r.FormValue("url")
	if url == "" {
		log.Println("blank url from", r.RemoteAddr)
		http.Error(w, "no url", 400)
		return
	}
	peerlist.announce <- url
	peersHandler(w, r, params)
}

func peersHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	group := params.ByName("group")
	if group == "" {
		http.Error(w, "no group", 400)
		return
	}
	peerlist := getPeers(group)

	recv := make(chan []string)
	peerlist.req <- recv
	list := <-recv
	data, _ := json.Marshal(list)
	w.Write(data)
}
