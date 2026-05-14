package handler

import (
	"go-chat/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	peerIDStr := c.Query("peer_id")
	peerID, err := strconv.ParseUint(peerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid peer_id"})
		return
	}

	messages, err := repository.GetMessages(userID, uint(peerID), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
