package handler

import (
	"encoding/json"
	"fmt"
	"go-chat/internal/chat"
	"go-chat/internal/model"
	"go-chat/internal/repository"
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
		UserID: uint(c.GetFloat64("user_id")),
		Name:   username,
		Send:   make(chan []byte, 256),
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
		var msg struct {
			Type string `json:"type"`
			To   uint   `json:"to"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		if msg.Type == "pm" && msg.To > 0 && msg.Text != "" {
			timestamp := time.Now().Format("15:04:05")
			formatted := fmt.Sprintf("[PM:%s:%s] %s", client.Name, timestamp, msg.Text)

			repository.SaveMessage(&model.Message{
				SenderID:   client.UserID,
				ReceiverID: msg.To,
				Content:    msg.Text,
			})

			h.hub.SendToUser(msg.To, []byte(formatted))
			h.hub.SendToUser(client.UserID, []byte(formatted))
		}

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
