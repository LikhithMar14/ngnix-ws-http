package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/ws-ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "WebSocket server is running", "port": 8081, "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "WebSocket Server is running on port 8081! Connect to /ws")
	})

	fmt.Println("WebSocket Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("WebSocket client connected")

	welcome := fmt.Sprintf("Welcome to WebSocket server! Connected at %s", time.Now().Format(time.RFC3339))
	conn.WriteMessage(websocket.TextMessage, []byte(welcome))

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		fmt.Printf("Received message: %s\n", message)
		
		response := fmt.Sprintf("Echo: %s (received at %s)", message, time.Now().Format(time.RFC3339))
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}

	fmt.Println("WebSocket client disconnected")
}