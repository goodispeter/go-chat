package repository

import (
	"go-chat/internal/database"
	"go-chat/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
