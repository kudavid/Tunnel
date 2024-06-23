package main

import (
	"github.com/hashicorp/yamux"
	"log"
	"net"
	"net/http"
)

func Listen(proxyAddr string) (net.Listener, error) {
	conn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		return nil, err
	}

	return yamux.Client(conn, nil)
}

func main() {
	l, err := Listen("127.0.0.1:8090")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World!\n"))
	})
	err = http.Serve(l, nil)
	log.Fatal(err)
}
