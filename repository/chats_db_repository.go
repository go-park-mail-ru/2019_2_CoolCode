package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
)

type ChatsDBRepository struct {
	db *sql.DB
}

func (c *ChatsDBRepository) GetWorkspaceByID(userID uint64) (models.Workspace, error) {
	panic("implement me")
}

func (c *ChatsDBRepository) GetWorkspaces() ([]models.Workspace, error) {
	panic("implement me")
}

func (c *ChatsDBRepository) PutWorkspace(workspace *models.Workspace) error {
	panic("implement me")
}

func (c *ChatsDBRepository) PutChannel(userID uint64, channel *models.Channel) error {
	panic("implement me")
}

func (c *ChatsDBRepository) UpdateWorkspace(workspaceID uint64, workspace *models.Workspace) error {
	panic("implement me")
}

func (c *ChatsDBRepository) GetChannelByID(channelID uint64) (models.Channel, error) {
	panic("implement me")
}

func (c *ChatsDBRepository) UpdateChannel(channelID uint64, channel *models.Channel) error {
	panic("implement me")
}

func (c *ChatsDBRepository) RemoveWorkspace(workspaceID uint64) error {
	panic("implement me")
}

func (c *ChatsDBRepository) RemoveChannel(channelID uint64) error {
	panic("implement me")
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
