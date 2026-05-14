package main

import (
	"go-chat/internal/chat"
	"go-chat/internal/config"
	"go-chat/internal/handler"
	"go-chat/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	hub := chat.NewHub()
	go hub.Run()

	ws := handler.NewWsHandler(hub)

	r := gin.Default()
	r.StaticFS("/web", http.Dir("./web"))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/index.html")
	})
	api := r.Group("/api")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	protected.GET("/ws", ws.HandleWebSocket)

	r.Run(":8080")
}
