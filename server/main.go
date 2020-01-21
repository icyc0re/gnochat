package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	host = flag.String("host", "localhost", "server hostname")
	port = flag.String("port", "5003", "port on which server listens to")
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection: %s", r.RemoteAddr)
}

func main() {
	// start server
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/echo", handleConnections)

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}