package chat

import (
	"fmt"
	"strings"
	"sync"
)

type Client struct {
	UserID uint
	Name   string
	Send   chan []byte
}

type Hub struct {
	Clients         map[*Client]bool
	ClientsByUserID map[uint]*Client
	Broadcast       chan []byte
	Register        chan *Client
	Unregister      chan *Client
	mu              sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:         make(map[*Client]bool),
		ClientsByUserID: make(map[uint]*Client),
		Broadcast:       make(chan []byte),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.ClientsByUserID[client.UserID] = client
			h.broadcastUserList()
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientsByUserID, client.UserID)
				close(client.Send)
			}
			h.broadcastUserList()
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			h.sendAll(message)
			h.mu.Unlock()
		}
	}
}

func (h *Hub) OnlineNames() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	names := make([]string, 0, len(h.Clients))
	for client := range h.Clients {
		names = append(names, client.Name)
	}
	return names
}

func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.ClientsByUserID[userID]; ok {
		select {
		case client.Send <- message:
		default:
		}
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

func (h *Hub) broadcastUserList() {
	parts := make([]string, 0, len(h.Clients))
	for client := range h.Clients {
		parts = append(parts, fmt.Sprintf("%d:%s", client.UserID, client.Name))
	}
	msg := fmt.Appendf(nil, "[System] Users: %s", strings.Join(parts, ","))
	h.sendAll(msg)
}
