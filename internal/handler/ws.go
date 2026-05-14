package handler

import (
	"fmt"
	"go-chat/internal/chat"
	"net/http"
	"time"

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
	username := c.GetString("username")
	if username == "" {
		username = "anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade failed:", err)
		return
	}

	client := &chat.Client{
		Name: username,
		Send: make(chan []byte, 256),
	}

	h.hub.Register <- client

	go h.readPump(conn, client)
	go h.writePump(conn, client)
}

func (h *WsHandler) readPump(conn *websocket.Conn, client *chat.Client) {
	defer func() {
		h.hub.Unregister <- client
	}()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fromatted := fmt.Sprintf("[%s]%s: %s", time.Now().Format("15:04:05"), client.Name, string(message))
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
