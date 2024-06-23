package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World!\n"))
	})
	err = http.Serve(l, nil)
	log.Fatal(err)
}
