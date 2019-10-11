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
	chats,err:=c.repository.GetChats()
	var userChats []models.Chat
	if err!=nil{
		return chats,err
	}
	for _,chat:=range chats{
		if contains(chat.Members,ID){
			userChats=append(userChats,chat)
		}
	}
	return userChats,nil
}

func (c chatsUseCase) GetChatByID(ID uint64) (models.Chat, error) {
	return c.repository.GetChatByID(ID)
}

func (c chatsUseCase) PutChat(Chat *models.Chat) error {
	return c.repository.PutChat(Chat)
}

func (c chatsUseCase) Contains(Chat models.Chat) error {
	return c.repository.Contains(Chat)
}

func NewChatsUseCase(repo repository.ChatsRepository) ChatsUseCase {
	return chatsUseCase{
		repository:repo,
	}
}

func contains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}