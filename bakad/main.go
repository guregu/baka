package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

var bind = flag.String("bind", ":1337", "bind address")
var timeout = flag.String("timeout", "1m", "dead peer timeout")
var peerlist *peers

func main() {
	flag.Parse()

	purgeTime, err := time.ParseDuration(*timeout)
	if err != nil {
		panic(err)
	}

	setup(purgeTime)

	log.Println("starting bakad:", *bind)
	http.ListenAndServe(*bind, nil)
}

func setup(purgeTime time.Duration) {
	peerlist = newPeers(purgeTime)

	http.HandleFunc("/peers", peersHandler)
	http.HandleFunc("/announce", announceHandler)
}

func announceHandler(w http.ResponseWriter, r *http.Request) {
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
	//w.Write([]byte("ok"))
	peersHandler(w, r)
}

func peersHandler(w http.ResponseWriter, r *http.Request) {
	recv := make(chan []string)
	peerlist.req <- recv
	list := <-recv
	data, _ := json.Marshal(list)
	w.Write(data)
}
