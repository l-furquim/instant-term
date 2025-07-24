package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	SERVER = "ws://localhost:9090/ws"
)

func main() {
	alias := flag.String("alias", "user", "Your alias in the chat")

	flag.Parse()

	conn, _, err := websocket.DefaultDialer.Dial(SERVER, nil)

	if err != nil {
		log.Fatalf("Error estabilishing the connection to the websocket server %s", err)
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte(*alias))
	if err != nil {
		log.Fatal("Error sending the alias:", err)
	}

	// Ctrl+C handler
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error while reading the message:", err)
				return
			}
			fmt.Printf("%s\n", msg)
		}
	}()

	fmt.Println("Type messages and press Enter to send (quit to exit the cli):")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
			break
		}
		input = strings.TrimSpace(input)

		if input == "quit" {
			signal.Notify(interrupt, os.Interrupt)
		}

		conn.WriteMessage(websocket.TextMessage, []byte(input))
	}

}
