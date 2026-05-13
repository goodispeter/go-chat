package handler

import (
	"fmt"
	"go-chat/internal/chat"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHandler struct {
	hub *chat.Hub
}

func NewWsHandler(hub *chat.Hub) *WsHandler {
	return &WsHandler{hub: hub}
}

func (h *WsHandler) HandleWebSocket(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade failed:", err)
		return
	}

	client := &chat.Client{
		Name: name,
		Send: make(chan []byte, 256),
	}

	h.hub.Register <- client

	h.hub.Broadcast <- []byte(fmt.Sprintf("[System] %s join chat room", name))

	go h.readPump(conn, client)
	go h.writePump(conn, client)
}

func (h *WsHandler) readPump(conn *websocket.Conn, client *chat.Client) {
	defer func() {
		h.hub.Unregister <- client
		h.hub.Broadcast <- []byte(fmt.Sprintf("[System] %s leave the chat room", client.Name))
	}()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fromatted := fmt.Sprintf("%s: %s", client.Name, string(message))
		h.hub.Broadcast <- []byte(fromatted)
	}
}

func (h *WsHandler) writePump(conn *websocket.Conn, client *chat.Client) {
	defer conn.Close()
	for message := range client.Send {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
