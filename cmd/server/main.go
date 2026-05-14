package main

import (
	"fmt"
	"go-chat/internal/chat"
	"go-chat/internal/config"
	"go-chat/internal/database"
	"go-chat/internal/handler"
	"go-chat/internal/middleware"
	"go-chat/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	if err := database.Connect(); err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	// Auto-create tables
	database.DB.AutoMigrate(&model.User{}, &model.Message{})

	hub := chat.NewHub()
	go hub.Run()

	ws := handler.NewWsHandler(hub)

	r := gin.Default()
	r.StaticFS("/web", http.Dir("./web"))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web/login.html")
	})
	api := r.Group("/api")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	protected.GET("/ws", ws.HandleWebSocket)
	protected.GET("/api/messages", handler.GetMessages)

	r.Run(":8080")
}
