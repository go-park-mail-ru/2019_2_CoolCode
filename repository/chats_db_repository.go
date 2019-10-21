package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
)

type ChatsDBRepository struct {
	db *sql.DB
}

func (c *ChatsDBRepository) GetChatByID(ID uint64) (models.Chat, error) {
	panic("implement me")
}

func (c *ChatsDBRepository) PutChat(Chat *models.Chat) error {
	panic("implement me")
}

func (c *ChatsDBRepository) Contains(Chat models.Chat) error {
	panic("implement me")
}

func (c *ChatsDBRepository) GetChats() ([]models.Chat, error) {
	panic("implement me")
}

func NewChatsDBRepository(db *sql.DB) ChatsRepository {
	return &ChatsDBRepository{db: db}
}
