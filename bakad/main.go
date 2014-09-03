package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var bind = flag.String("bind", ":1337", "bind address")
var peerlist *peers

func main() {
	flag.Parse()

	peerlist = newPeers()

	http.HandleFunc("/", test)
	http.HandleFunc("/announce", announce)

	log.Println("starting bakad:", *bind)
	http.ListenAndServe(*bind, nil)
}

func announce(w http.ResponseWriter, r *http.Request) {
	peerlist.announce <- r.RemoteAddr
	w.Write([]byte("ok"))
}

func test(w http.ResponseWriter, r *http.Request) {
	recv := make(chan []string)
	peerlist.req <- recv
	list := <-recv
	data, _ := json.Marshal(list)
	w.Write(data)
}
