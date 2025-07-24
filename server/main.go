package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name string
	Conn *websocket.Conn
	Send chan []byte
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan []byte)
	upgrader  = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Websocket server running in :9090/ws")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error creating the upgrade:", err)
		return
	}

	_, nameBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading the name:", err)
		return
	}
	name := string(nameBytes)

	client := &Client{
		Name: name,
		Conn: conn,
		Send: make(chan []byte),
	}

	clients[client] = true

	log.Printf("A connection is made from the addr %s \n", conn.RemoteAddr().String())

	joinMsg := fmt.Sprintf("%s joined the chat!", client.Name)
	broadcast <- []byte(joinMsg)

	go handleClientRead(client)
	go handleClientWrite(client)
}

func handleClientRead(client *Client) {
	defer func() {
		delete(clients, client)
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}
		formatted := fmt.Sprintf("[%s] %s", client.Name, msg)
		broadcast <- []byte(formatted)
	}
}

func handleClientWrite(client *Client) {
	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error sending:", err)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}
