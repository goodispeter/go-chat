package chat

import (
	"fmt"
	"sync"
)

type Client struct {
	Name string
	Send chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) sendAll(message []byte) {
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.sendAll(fmt.Appendf(nil, "[System] %s join chat room", client.Name))
			h.sendAll(fmt.Appendf(nil, "[System] Online: %d", len(h.Clients)))
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.sendAll(fmt.Appendf(nil, "[System] %s leave the chat room", client.Name))
			h.sendAll(fmt.Appendf(nil, "[System] Online: %d", len(h.Clients)))
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			h.sendAll(message)
			h.mu.Unlock()
		}
	}
}
