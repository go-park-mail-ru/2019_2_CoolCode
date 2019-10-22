package repository

import "github.com/go-park-mail-ru/2019_2_CoolCode/models"

type ChatsRepository interface {
	GetChatByID(ID uint64) (models.Chat, error)
	PutChat(Chat *models.Chat) error
	Contains(Chat models.Chat) error
	GetChats(userID uint64) ([]models.Chat, error)
	GetWorkspaceByID(userID uint64) (models.Workspace, error)
	GetWorkspaces(userID uint64) ([]models.Workspace, error)
	PutWorkspace(workspace *models.Workspace) error
	PutChannel(channel *models.Channel) error
	UpdateWorkspace(workspace *models.Workspace) error
	GetChannelByID(channelID uint64) (models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	RemoveWorkspace(workspaceID uint64) error
	RemoveChannel(channelID uint64) error
	RemoveChat(chatID uint64) error
}
