package hub

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*websocket.Conn]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *websocket.Conn

	// Unregister requests from clients.
	unregister chan *websocket.Conn
}

func NewHub() *Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
	go h.run()
	return h
}

func (h *Hub) Brodcast(msg string) error {
	select {
	case h.broadcast <- []byte(msg):
		return nil
	case <-time.After(1 * time.Second):
		return fmt.Errorf("get time out when broadcast ws message %s", msg)
	}
}

func (h *Hub) Register(con *websocket.Conn) error {
	select {
	case h.register <- con:
		return nil
	case <-time.After(1 * time.Second):
		return fmt.Errorf("get time out when register ws %s", con)
	}
}

func (h *Hub) UnRegister(con *websocket.Conn) error {
	select {
	case h.unregister <- con:
		return nil
	case <-time.After(1 * time.Second):
		return fmt.Errorf("get time out when register ws %s", con)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			log.Println("get brodcast message, brodcasting", message)
			for client := range h.clients {
				if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Printf("get error when send ws data, error: %s\n", err)
					h.UnRegister(client)
				}
			}
		}
	}

}
