package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	host = flag.String("host", "localhost", "server hostname")
	port = flag.String("port", "5003", "port on which server listens to")
)

func connect() *websocket.Conn {
	addr := fmt.Sprintf("%s:%s", *host, *port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	fmt.Printf("connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("ERROR dial:",err)
		os.Exit(1)
	}

	return conn
}

func main() {
	flag.Parse()

	conn := connect()
	defer conn.Close()

	for i := 0; i < 20; i++ {
		msg := []byte(fmt.Sprintf("msg %d", i))
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("ERROR write:", err)
		}

		time.Sleep(time.Millisecond * 100)
	}
}