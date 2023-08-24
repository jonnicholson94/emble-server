package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Websocket(w http.ResponseWriter, r *http.Request) {

	// Initialise the connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	fmt.Println("Client connected")

	for {
		message, p, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Println(message)

		query := string(p)

		fmt.Println("Received query:", query)

	}
}
