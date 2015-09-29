package main

import (
	"log"
	"fmt"
	"net/http"
	"net/url"
)

func commandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	if cmd == "BLOCK" {
	}
	switch cmd {
	case "BLOCK" :
		target, _ = url.Parse(shimUrl)
	case "ALLOW" :
		target, _ = url.Parse(*targetUrl)
	default: log.Println("Unkonwn command")
	}
}

func StartCommander(host string, port int) {

	addr := fmt.Sprintf("%s:%d", host, port)
	svr := http.NewServeMux()
	svr.HandleFunc("/", commandHandler)

	log.Printf("Starting http command interface at %s", addr)

	err := http.ListenAndServe(addr, svr)

	if err != nil {
		log.Fatal("ERROR: failed to start command interface", err)

	}
}

var shimUrl string

func StartShimServer(host string, port int) {

	addr := fmt.Sprintf("%s:%d", host, port)
	shimUrl = fmt.Sprintf("http://%s:%d", host, port)

	svr := http.NewServeMux()
	svr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nothing to see here")
	})

	log.Printf("Starting http command interface at %s", addr)

	err := http.ListenAndServe(addr, svr)

	if err != nil {
		log.Fatal("ERROR: failed to start shim server", err)

	}
}
