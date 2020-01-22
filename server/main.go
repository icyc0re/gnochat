package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	host = flag.String("host", "localhost", "server hostname")
	port = flag.String("port", "5003", "port on which server listens to")
	upgrader = websocket.Upgrader{} // default configuration
	clients = make(map[*websocket.Conn]string)
)

func closeConnection(conn *websocket.Conn) {
	conn.Close()

	log.Printf("Disconnect: %s", clients[conn])
	delete(clients, conn)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection: %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR upgrade:", err)
		return
	}

	defer closeConnection(conn)

	clients[conn] = "username"

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Disconnect: %s", r.RemoteAddr)
				conn.WriteMessage(mt, msg)
				return
			} else if websocket.IsUnexpectedCloseError(err) {
				log.Printf("Disconnect: %s", r.RemoteAddr)
				return
			} else {
				log.Println("ERROR read:", err)
				break
			}
		}

		log.Printf("recv: %s", msg)
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("ERROR write:", err)
			break
		}
	}
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