package service

import (
	"go-chat/internal/model"
	"go-chat/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return repository.CreateUser(user)
}

func Login(username, password string) error {
	user, err := repository.FindByUsername(username)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
