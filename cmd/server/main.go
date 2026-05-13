package main

import (
	"go-chat/internal/chat"
	"go-chat/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := chat.NewHub()
	go hub.Run()

	ws := handler.NewWsHandler(hub)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "go-chat is running")
	})
	r.GET("/ws", ws.HandleWebSocket)

	r.Run(":8080")
}
