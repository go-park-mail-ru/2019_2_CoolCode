package useCase

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

var messageUseCase = MessageUseCaseImpl{
	repository: &repository.MessageRepositoryMock{},
	chats:      &ChatsUseCaseMock{},
}

func TestMessageUseCaseImpl_GetChatMessages(t *testing.T) {
	chatID := 0
	messageID := 0

	//test internal error
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return false, errors.New("Internal error")
		},
	}

	_, err := messageUseCase.GetChatMessages(uint64(chatID), uint64(messageID))

	assert.NotNil(t, err)

	//test not permission
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return false, nil
		},
	}

	_, err = messageUseCase.GetChatMessages(uint64(chatID), uint64(messageID))

	assert.NotNil(t, err)

	//test success
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return true, nil
		},
	}
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessagesByChatIDFunc: func(chatID uint64) (messages models.Messages, e error) {
			return models.Messages{}, nil
		},
	}

	_, err = messageUseCase.GetChatMessages(uint64(chatID), uint64(messageID))

	assert.Nil(t, err)
}

func TestMessageUseCaseImpl_SaveMessage(t *testing.T) {
	testMessage := &models.Message{
		ID:     0,
		ChatID: 0,
	}

	//test internal error
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return false, errors.New("Internal error")
		},
	}

	_, err := messageUseCase.SaveChatMessage(testMessage)

	assert.NotNil(t, err)

	//test not permission
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return false, nil
		},
	}

	_, err = messageUseCase.SaveChatMessage(testMessage)

	assert.NotNil(t, err)

	//test success
	messageUseCase.chats = &ChatsUseCaseMock{
		CheckChatPermissionFunc: func(userID uint64, chatID uint64) (b bool, e error) {
			return true, nil
		},
	}
	messageUseCase.repository = &repository.MessageRepositoryMock{
		PutMessageFunc: func(message *models.Message) (u uint64, e error) {
			return 0, nil
		},
	}

	_, err = messageUseCase.SaveChatMessage(testMessage)

	assert.Nil(t, err)
}

func TestMessageUseCaseImpl_EditMessage(t *testing.T) {
	userID := 0
	testMessage := &models.Message{
		ID:     0,
		ChatID: 0,
	}

	//test internal error
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{}, errors.New("Internal error")
		},
	}

	err := messageUseCase.EditMessage(testMessage, uint64(userID))

	assert.NotNil(t, err)

	//test not permission
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{AuthorID: 1}, errors.New("Internal error")
		},
	}

	err = messageUseCase.EditMessage(testMessage, uint64(userID))

	assert.NotNil(t, err)

}

func TestMessageUseCaseImpl_DeleteMessage(t *testing.T) {
	authorID := 0
	messageID := 0

	//test internal error
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{}, errors.New("Internal error")
		},
	}

	err := messageUseCase.DeleteMessage(uint64(messageID), uint64(authorID))

	assert.NotNil(t, err)

	//test not permission
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{AuthorID: 1}, nil
		},
	}

	err = messageUseCase.DeleteMessage(uint64(messageID), uint64(authorID))

	assert.NotNil(t, err)

	//test success
	authorID = 1
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{AuthorID: 1}, nil
		},
		RemoveMessageFunc: func(messageID uint64) error {
			return nil
		},
	}

	err = messageUseCase.DeleteMessage(uint64(messageID), uint64(authorID))

	assert.Nil(t, err)
}

func TestMessageUseCaseImpl_HideMessageForAuthor(t *testing.T) {
	authorID := 0
	messageID := 0

	//test internal error
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{}, errors.New("Internal error")
		},
	}

	err := messageUseCase.HideMessageForAuthor(uint64(messageID), uint64(authorID))

	assert.NotNil(t, err)

	//test not permission
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{AuthorID: 1}, nil
		},
	}

	err = messageUseCase.HideMessageForAuthor(uint64(messageID), uint64(authorID))

	assert.NotNil(t, err)

	//test success
	authorID = 1
	messageUseCase.repository = &repository.MessageRepositoryMock{
		GetMessageByIDFunc: func(messageID uint64) (message *models.Message, e error) {
			return &models.Message{AuthorID: 1}, nil
		},
		HideMessageForAuthorFunc: func(userID uint64) error {
			return nil
		},
	}

	err = messageUseCase.HideMessageForAuthor(uint64(messageID), uint64(authorID))

	assert.Nil(t, err)
}
