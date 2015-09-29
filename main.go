package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
	"strings"
	"flag"
	"fmt"
	"sync"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func switchTarget(req *http.Request) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

var (
	targetLock  sync.Mutex
	target *url.URL
)

func saveTarget(u string) {
	targetLock.Lock()
	target, _ = url.Parse(u)
	targetLock.Unlock()
}

var targetUrl = flag.String("t", "http://localhost:3000", "URL of target site")
var host = flag.String("b", "0.0.0.0", "Binds to the specified IP")
var port = flag.Int("p", 8080, "Runs on the specified port")

func main() {
	flag.Parse()

	saveTarget(*targetUrl)

	go StartCommander(*host, 9090)
	go StartShimServer(*host, 9091)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = switchTarget
	http.Handle("/", proxy)

	log.Printf("Starting proxy interface at %s", addr)
	err := http.ListenAndServe( addr, nil)

	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

