package useCase

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"net/http"
)

//go:generate moq -out chats_ucase_mock.go . ChatsUseCase

type ChatsUseCase interface {
	CheckChatPermission(userID uint64, chatID uint64) (bool, error)
	GetChatByID(userID uint64, ID uint64) (models.Chat, error)
	GetChatsByUserID(ID uint64) ([]models.Chat, error)
	PutChat(Chat *models.Chat) (uint64, error)
	Contains(Chat models.Chat) error
	GetWorkspaceByID(userID uint64, ID uint64) (models.Workspace, error)
	GetWorkspacesByUserID(ID uint64) ([]models.Workspace, error)
	CreateWorkspace(room *models.Workspace) (uint64, error)
	CreateChannel(channel *models.Channel) (uint64, error)
	GetChannelByID(userID uint64, ID uint64) (models.Channel, error)
	EditWorkspace(userID uint64, room *models.Workspace) error
	EditChannel(userID uint64, channel *models.Channel) error
	LogoutFromWorkspace(userID uint64, workspaceID uint64) error
	LogoutFromChannel(userID uint64, channelID uint64) error
	DeleteWorkspace(userID uint64, workspaceID uint64) error
	DeleteChannel(userID uint64, channelID uint64) error
	DeleteChat(userID uint64, chatId uint64) error
}

type ChatsUseCaseImpl struct {
	repository      repository.ChatsRepository
	usersRepository repository.UserRepo
}

func (c *ChatsUseCaseImpl) PutChat(Chat *models.Chat) (uint64, error) {
	return c.repository.PutChat(Chat)
}

func (c *ChatsUseCaseImpl) GetChatByID(userID uint64, ID uint64) (models.Chat, error) {
	chat, err := c.repository.GetChatByID(ID)
	if err != nil {
		return chat, err
	}
	if !contains(chat.Members, userID) {
		return chat, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	lastMessage, err := c.repository.GetMessageLast(ID)
	if err != nil {
		fmt.Println("BLYA YPALI", err)
		return chat, nil
	}
	chat.LastMessage = lastMessage.Text
	return chat, nil
}

func (c *ChatsUseCaseImpl) GetChatsByUserID(ID uint64) ([]models.Chat, error) {
	chats, err := c.repository.GetChats(ID)
	var userChats []models.Chat
	if err != nil {
		return chats, err
	}
	for _, chat := range chats {
		var memberID uint64
		if contains(chat.Members, ID) {
			for _, userID := range chat.Members {
				if userID != ID {
					memberID = userID
					break
				}
			} //get chat name
			user, _ := c.usersRepository.GetUserByID(memberID)
			if user.Username != "" {
				chat.Name = user.Username
			}
			userChats = append(userChats, chat)
		}
	}
	return userChats, nil
}

func (c *ChatsUseCaseImpl) DeleteChat(userID uint64, chatID uint64) error {
	ok, err := c.CheckChatPermission(userID, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	_, err = c.repository.RemoveChat(chatID)
	return err
}

func (c *ChatsUseCaseImpl) CreateWorkspace(workspace *models.Workspace) (uint64, error) {
	return c.repository.PutWorkspace(workspace)
}

func (c *ChatsUseCaseImpl) GetWorkspaceByID(userID uint64, ID uint64) (models.Workspace, error) {
	workspace, err := c.repository.GetWorkspaceByID(ID)
	if err != nil {
		return workspace, err
	}
	if !contains(workspace.Members, userID) {
		return workspace, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return workspace, nil
}

func (c *ChatsUseCaseImpl) GetWorkspacesByUserID(ID uint64) ([]models.Workspace, error) {
	workspaces, err := c.repository.GetWorkspaces(ID)
	var userWorkspaces []models.Workspace
	if err != nil {
		return workspaces, err
	}
	for _, workspace := range workspaces {
		if contains(workspace.Members, ID) {
			userWorkspaces = append(userWorkspaces, workspace)
		}
	}
	return userWorkspaces, nil
}

func (c *ChatsUseCaseImpl) EditWorkspace(userID uint64, workspace *models.Workspace) error {
	editWorkspace, err := c.repository.GetWorkspaceByID(workspace.ID)
	if err != nil {
		return err
	}
	if !contains(editWorkspace.Admins, userID) {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	workspace.Channels = editWorkspace.Channels
	workspace.CreatorID = editWorkspace.CreatorID
	return c.repository.UpdateWorkspace(workspace)
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) DeleteWorkspace(userID uint64, workspaceID uint64) error {
	deleting, err := c.repository.GetWorkspaceByID(workspaceID)
	if err != nil {
		return err
	}
	if userID != deleting.CreatorID {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	_, err = c.repository.RemoveWorkspace(workspaceID)
	return err
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) LogoutFromWorkspace(userID uint64, workspaceID uint64) error {
	editWorkspace, err := c.repository.GetWorkspaceByID(workspaceID)
	if err != nil {
		return err
	}
	if !contains(editWorkspace.Members, userID) {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	editWorkspace.Members = removeElement(editWorkspace.Members, userID)
	return c.repository.UpdateWorkspace(&editWorkspace)
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) CreateChannel(channel *models.Channel) (uint64, error) {
	workspace, err := c.repository.GetWorkspaceByID(channel.WorkspaceID)
	if err != nil {
		return 0, err
	}
	if !contains(workspace.Admins, channel.CreatorID) {
		return 0, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return c.repository.PutChannel(channel)
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) GetChannelByID(userID, ID uint64) (models.Channel, error) {
	channel, err := c.repository.GetChannelByID(ID)
	if err != nil {
		return channel, err
	}
	if !contains(channel.Members, userID) {
		return channel, models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	return channel, nil
}

func (c *ChatsUseCaseImpl) EditChannel(userID uint64, channel *models.Channel) error {
	editChannel, err := c.repository.GetChannelByID(channel.ID)
	if err != nil {
		return err
	}
	if !contains(editChannel.Admins, userID) {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	channel.TotalMSGCount = editChannel.TotalMSGCount
	channel.CreatorID = editChannel.CreatorID
	return c.repository.UpdateChannel(channel)
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) LogoutFromChannel(userID uint64, channelID uint64) error {
	editChannel, err := c.repository.GetChannelByID(channelID)
	if err != nil {
		return err
	}
	if !contains(editChannel.Members, userID) {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	editChannel.Members = removeElement(editChannel.Members, userID)
	return c.repository.UpdateChannel(&editChannel)
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) DeleteChannel(userID uint64, channelID uint64) error {
	deletingRoom, err := c.repository.GetChannelByID(channelID)
	if err != nil {
		return err
	}
	if userID != deletingRoom.CreatorID {
		return models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:(")
	}
	_, err = c.repository.RemoveChannel(channelID)
	return err
	//TODO: отправить уведомление всем открытм ws
}

func (c *ChatsUseCaseImpl) CheckChatPermission(userID uint64, chatID uint64) (bool, error) {
	_, err := c.GetChatByID(userID, chatID)
	return err == nil, err //TODO: плохо
}

func (c *ChatsUseCaseImpl) Contains(Chat models.Chat) error {
	return c.repository.Contains(Chat)
}

func NewChatsUseCase(repo repository.ChatsRepository, usersRepo repository.UserRepo) ChatsUseCase {
	return &ChatsUseCaseImpl{
		repository:      repo,
		usersRepository: usersRepo,
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

func removeElement(s []uint64, e uint64) []uint64 {
	var index int
	for i, elem := range s {
		if elem == e {
			index = i
		}
	}
	return append(s[:index], s[index+1:]...)
}
