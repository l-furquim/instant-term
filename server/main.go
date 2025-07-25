package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

const (
	QUIT_COMMAND  = "/quit"
	HELP_COMMAND  = "/help"
	WISP_COMMAND  = "/w"
	COMMAND_LISTS = "\033[34m/quit -> Quit the server and cli \n/help -> Show the commands list \n/w subject_name message\033[0m"
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
		disconnectClient(client)
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}

		command := strings.TrimSpace(strings.ToLower(string(msg)))

		switch {
		case strings.Contains(command, WISP_COMMAND):
			handleWispMessage(client, &msg)

		case command == HELP_COMMAND:
			client.Send <- []byte(COMMAND_LISTS)

		case command == QUIT_COMMAND:
			disconnectClient(client)
			return

		default:
			formatted := fmt.Sprintf("\033[35m[%s]\033[0m %s", client.Name, msg)
			broadcast <- []byte(formatted)
		}
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
			ms := string(msg)

			if strings.Contains(ms, client.Name) {
				ms = strings.Replace(ms, "\033[35m", "\033[32m", 1)
			}

			select {
			case client.Send <- []byte(ms):
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

func disconnectClient(client *Client) {
	delete(clients, client)
	client.Conn.Close()

	log.Printf("Connection closed from addr %s \n", client.Conn.RemoteAddr().String())
}

func handleWispMessage(client *Client, msg *[]byte) {
	log.Println("Not implemented yeat.")
}
