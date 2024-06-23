package main

import (
	"github.com/hashicorp/yamux"
	"io"
	"log"
	"net"
	"net/http"
)

func Listen(proxyAddr string) (net.Listener, error) {
	var conn io.ReadWriteCloser
	// TODO create conn from HTTP/2
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
