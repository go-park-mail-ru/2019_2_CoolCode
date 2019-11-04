package repository

import "github.com/go-park-mail-ru/2019_2_CoolCode/models"

//go:generate moq -out message_repo_mock.go . MessageRepository

type MessageRepository interface {
	PutMessage(message *models.Message) (uint64, error)
	GetMessageByID(messageID uint64) (*models.Message, error)
	GetMessagesByChatID(chatID uint64) (models.Messages, error)
	RemoveMessage(messageID uint64) error
	UpdateMessage(message *models.Message) error
	HideMessageForAuthor(userID uint64) error
}
