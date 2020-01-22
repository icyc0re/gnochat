package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type UUID int64

type Message struct {
	uuid UUID
	msg []byte
}

func (m *Message) Prepare() []byte {
	return []byte(fmt.Sprintf("%d:%s:%s", m.uuid, clients[m.uuid].username, string(m.msg)))
}

type UserData struct {
	conn *websocket.Conn
	username string
}

var (
	host = flag.String("host", "localhost", "server hostname")
	port = flag.String("port", "5003", "port on which server listens to")
	upgrader = websocket.Upgrader{} // default configuration
	clients = make(map[UUID]UserData)
	broadcastChan = make(chan Message, 100)
	maxId struct {
		currentValue UUID
		mux sync.Mutex
	}
)

func closeConnection(uuid UUID) {
	user := clients[uuid]
	user.conn.Close()

	log.Printf("Disconnect: %s", user.username)
	delete(clients, uuid)
}

func validUsername(username string) bool {
	return !strings.Contains(username, ":")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("New connection: %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR upgrade:", err)
		return
	}

	maxId.mux.Lock()
	maxId.currentValue++
	userId := maxId.currentValue
	maxId.mux.Unlock()

	defer closeConnection(userId)

	// expect username as first message
	mt, msg, err := conn.ReadMessage()
	if err != nil || mt != websocket.TextMessage {
		log.Println("ERROR, expected username as first message after connection!")
		return
	}
	
	// check validity of username
	username := string(msg)
	if !validUsername(username) {
		log.Println("ERROR: invalid username!")
		return
	}

	// answer with userId
	err = conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(userId)))
	if err != nil {
		log.Printf("ERROR sending userId to %s\n", username)
		return
	}
	
	clients[userId] = UserData{ conn, username }

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				conn.WriteMessage(mt, msg)
				return
			} else if websocket.IsUnexpectedCloseError(err) {
				return
			} else {
				log.Println("ERROR read:", err)
				break
			}
		}

		broadcastChan <- Message{ userId, msg }
	}
}

func broadcastMessages() {
	for {
		msg := <-broadcastChan

		for _, clientData := range clients {
			err := clientData.conn.WriteMessage(websocket.TextMessage, msg.Prepare())
			if err != nil {
				log.Println("ERROR write:", err)
				break
			}
		}
	}
}

func main() {
	// start server
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/echo", handleConnections)

	go broadcastMessages()

	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
