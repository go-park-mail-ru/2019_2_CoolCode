package repository

import "github.com/go-park-mail-ru/2019_2_CoolCode/models"

type ChatsRepository interface {
	GetChatByID(ID uint64)  (models.Chat,error)
	PutChat(Chat *models.Chat) error
	Contains(Chat models.Chat) error
	GetChats() ([]models.Chat,error)
}
