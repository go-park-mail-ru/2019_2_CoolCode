package useCase

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
)

type ChatsUseCase interface {
	GetChatByID(ID uint64)  (models.Chat,error)
	GetChatByUserID(ID uint64)  ([]models.Chat,error)
	PutChat(Chat *models.Chat) error
	Contains(Chat models.Chat) error
}

type chatsUseCase struct {
	repository repository.ChatsRepository
}

func (c chatsUseCase) GetChatByUserID(ID uint64) ([]models.Chat, error) {
	panic("implement me")
}

func (c chatsUseCase) GetChatByID(ID uint64) (models.Chat, error) {
	panic("implement me")
}

func (c chatsUseCase) PutChat(Chat *models.Chat) error {
	panic("implement me")
}

func (c chatsUseCase) Contains(Chat models.Chat) error {
	panic("implement me")
}

func NewChatsUseCase(repo repository.ChatsRepository) ChatsUseCase {
	return chatsUseCase{
		repository:repo,
	}
}