package main

import (
	"context"
	"github.com/hashicorp/yamux"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type reverseProxy struct {
	session *yamux.Session
	m       sync.Mutex
}

func (rp *reverseProxy) NewConn(conn io.ReadWriteCloser) {
	session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	rp.m.Lock()
	rp.session = session
	rp.m.Unlock()
}

func (rp *reverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var session *yamux.Session
	rp.m.Lock()
	session = rp.session
	rp.m.Unlock()

	if session == nil {
		http.Error(w, "Session Closed", http.StatusInternalServerError)
		return
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return session.Open()
		},
	}
	httpRP := &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(r *httputil.ProxyRequest) {
			target := &url.URL{Scheme: "http", Host: "yamux", Path: "/"}
			r.SetURL(target)
		},
	}
	httpRP.ServeHTTP(w, r)
}

func main() {
	rp := &reverseProxy{}

	http.Handle("/app/", rp)
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser
		// TODO createConn
		rp.NewConn(conn)
	})
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}
