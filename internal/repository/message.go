package repository

import (
	"go-chat/internal/database"
	"go-chat/internal/model"
)

func SaveMessage(msg *model.Message) error {
	return database.DB.Create(msg).Error
}

func GetMessages(userID, peerID uint, limit int) ([]model.Message, error) {
	var messages []model.Message
	err := database.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, peerID, peerID, userID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
