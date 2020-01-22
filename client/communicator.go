package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type UUID int64

func ValidUUID(id UUID) bool {
	return id > 0
}

func Connect(host, port string) *websocket.Conn {
	addr := fmt.Sprintf("%s:%s", host, port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	log.Printf("Connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("ERROR dial:",err)
		return nil
	}

	return conn
}

func InitialHandshake(conn *websocket.Conn, username string) UUID {

	if err := conn.WriteMessage(websocket.TextMessage, []byte(username)); err != nil {
		log.Println("handshake write:", err)
		return -1
	}

	if _, msg, err := conn.ReadMessage(); err != nil {
		log.Println("handshake read:", err)
	} else if userId, err := strconv.ParseInt(string(msg), 10, 64); err != nil {
		log.Println("ERROR received invalid/corrupted user id from server")
	} else {
		return UUID(userId)
	}

	return -1
}

func Disconnect(conn *websocket.Conn) {

	done := make(chan bool)

	go func() {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye!"))
		if err != nil {
			log.Println("ERROR close:", err)
			done <- false
			return
		}

		// wait for close message
		for {
			_, _, err := conn.ReadMessage()
			if err != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Connection closed safely!")
				done <- true
				return
			}
		}
	} ()

	select {
	case <-done:
	case <-time.After(DISCONNECT_TIMEOUT):
	}
}

// Send messages from channel ch, over connection conn
func MessageSender(conn *websocket.Conn, ch chan string) {
	for {
		msg := []byte(<-ch)
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// Receive messages on connection conn, redirects to channel ch
func MessageReceiver(conn *websocket.Conn, ch chan string) {
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ERROR read:", err)
			break
		}

		if mt == websocket.TextMessage {
			ch<- string(msg)
		}
	}
}
