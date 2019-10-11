package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"sync"
)

type ChatArrayRepository struct {
	chats  map[uint64]*models.Chat
	mutex  sync.Mutex
	nextID uint64
}

func NewChatArrayRepository() ChatsRepository{
	return ChatArrayRepository{
		chats:  make(map[uint64]*models.Chat, 0),
		mutex:  sync.Mutex{},
		nextID: 0,
	}
}

func (c ChatArrayRepository) GetChatByID(ID uint64) (models.Chat, error) {
	var resultChat models.Chat
	c.mutex.Lock()
	if user, ok := c.chats[ID]; ok {
		return *user, nil
	}
	c.mutex.Unlock()
	return resultChat, errors.New("user not contains")
}

func (c ChatArrayRepository) PutChat(Chat *models.Chat) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if Chat.ID==0 {
		c.nextID++
		Chat.ID = c.nextID
	}
	c.chats[Chat.ID] = Chat

	return nil
}

func (c ChatArrayRepository) Contains(Chat models.Chat) error {
	//for _, v := range c.chats {
	//	if v. == v.Email {
	//		return true
	//	}
	//}
	return nil
}

func (c ChatArrayRepository) GetChats() ([]models.Chat,error) {
	c.mutex.Lock()
	var chatsSlice []models.Chat
	for _, chat := range c.chats {
		chatsSlice = append(chatsSlice, *chat)
	}
	return chatsSlice,nil
}


