package useCase

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"net/http"
)

//go:generate moq -out messages_ucase_mock.go . MessagesUseCase
type MessagesUseCase interface {
	SaveChatMessage(message *models.Message) (uint64, error)
	EditMessage(message *models.Message, userID uint64) error
	DeleteMessage(messageID uint64, userID uint64) error
	GetChatMessages(chatID uint64, userID uint64) (models.Messages, error)
	GetMessageByID(messageID uint64) (*models.Message, error)
	HideMessageForAuthor(messageID uint64, userID uint64) error
	SaveChannelMessage(message *models.Message) (uint64, error)
	GetChannelMessages(channelID uint64, userID uint64) (models.Messages, error)
	FindMessages(findString string, ID uint64) (models.Messages, error)
}

type MessageUseCaseImpl struct {
	repository repository.MessageRepository
	chats      ChatsUseCase
}

func NewMessageUseCase(repository repository.MessageRepository, chats ChatsUseCase) MessagesUseCase {
	return &MessageUseCaseImpl{
		repository: repository,
		chats:      chats,
	}
}

func (m *MessageUseCaseImpl) GetChatMessages(chatID uint64, userID uint64) (models.Messages, error) {
	permissionOk, err := m.chats.CheckChatPermission(userID, chatID)
	if err != nil {
		return models.Messages{}, err
	}
	if !permissionOk {
		return models.Messages{}, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}

	return m.repository.GetMessagesByChatID(chatID)
}

func (m *MessageUseCaseImpl) GetChannelMessages(chatID uint64, userID uint64) (models.Messages, error) {
	permissionOk, err := m.chats.CheckChannelPermission(userID, chatID)
	if err != nil {
		return models.Messages{}, err
	}
	if !permissionOk {
		return models.Messages{}, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}

	return m.repository.GetMessagesByChatID(chatID)
}

func (m *MessageUseCaseImpl) GetMessageByID(messageID uint64) (*models.Message, error) {
	return m.repository.GetMessageByID(messageID)
}

func (m *MessageUseCaseImpl) SaveChatMessage(message *models.Message) (uint64, error) {
	permissionOk, err := m.chats.CheckChatPermission(message.AuthorID, message.ChatID)
	if err != nil {
		return 0, err
	}
	if !permissionOk {
		return 0, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return m.repository.PutMessage(message)
}

func (m *MessageUseCaseImpl) SaveChannelMessage(message *models.Message) (uint64, error) {
	permissionOk, err := m.chats.CheckChannelPermission(message.AuthorID, message.ChatID)
	if err != nil {
		return 0, err
	}
	if !permissionOk {
		return 0, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return m.repository.PutMessage(message)
}

func (m *MessageUseCaseImpl) EditMessage(message *models.Message, userID uint64) error {
	DBmessage, err := m.repository.GetMessageByID(message.ID)
	if err != nil {
		return err
	}
	if userID != DBmessage.AuthorID {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return m.repository.UpdateMessage(message)
}

func (m *MessageUseCaseImpl) DeleteMessage(messageID uint64, userID uint64) error {
	message, err := m.repository.GetMessageByID(messageID)
	if err != nil {
		return err
	}
	if userID != message.AuthorID {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return m.repository.RemoveMessage(messageID)
}

func (m *MessageUseCaseImpl) HideMessageForAuthor(messageID uint64, userID uint64) error {
	message, err := m.repository.GetMessageByID(messageID)
	if err != nil {
		return err
	}
	if userID != message.AuthorID {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return m.repository.HideMessageForAuthor(messageID)
}

func (m *MessageUseCaseImpl) FindMessages(findString string, ID uint64) (models.Messages, error) {
	messages, err := m.repository.FindMessages(findString)
	if err != nil {
		return messages, err
	}
	result := models.Messages{}

	for _, message := range messages.Messages {
		ok, err := m.chats.CheckChatPermission(ID, message.ChatID)
		if err != nil {
			ok, _ := m.chats.CheckChannelPermission(ID, message.ChatID)
			if ok {
				result.Messages = append(result.Messages, message)
			}
			continue
		}
		if ok {
			result.Messages = append(result.Messages, message)
		}
	}
	return result, nil
}
