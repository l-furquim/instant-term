package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	content   chan []byte
	timestamp time.Time
}

type Client struct {
	ID    string
	Conn  *websocket.Conn
	alias string
	Send  chan []byte
}

const (
	PORT           = "9090"
	maxMessageSize = 1024
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
	clients = []Client{}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var name = flag.String("name", "user", "your alias in the room")

func server(w http.ResponseWriter, r *http.Request, name *string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error while receiving the message %s", err)
	}

	client := &Client{
		alias: *name,
		Conn:  conn,
		ID:    time.Now().String(),
		Send:  make(chan []byte, 1024),
	}

	clients = append(clients, *client)

	go handleWrite(client)
	go handleRead(client)
}

func main() {
	flag.Parse()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server(w, r, name)
	})

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("Error listening: ", err)
	} else {
		log.Printf("Server running..")
	}
}

func handleSendMessages(message *[]byte) {
	for n := range clients {
		select {
		case clients[n].Send <- *message:
		default:
			close(clients[n].Send)
		}
	}
}

func handleWrite(client *Client) {
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatalf("Error while writing the message %s", err)
				return
			}

			w.Write(message)

			n := len(client.Send)
			for range n {
				w.Write(newLine)
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				log.Fatalf("Error while closing the connection %s", err)
				return
			}
		}
	}
}

func handleRead(client *Client) {
	defer client.Conn.Close()

	client.Conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		aliasColored := []byte("\033[35m" + client.alias + "\033[0m: ")
		cleanMessage := bytes.Replace(message, newLine, space, -1)
		message = append(aliasColored, cleanMessage...)

		handleSendMessages(&message)
	}
}
