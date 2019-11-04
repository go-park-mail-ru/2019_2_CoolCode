package useCase

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var chatsUseCase = ChatsUseCaseImpl{}

func TestChatsUseCaseImpl_GetChatByID(t *testing.T) {
	chatID := 0
	userID := 0
	//Internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.GetChatByID(uint64(userID), uint64(chatID))
	assert.NotNil(t, err)

	//Not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{Members: []uint64{1}}, nil
		},
	}

	_, err = chatsUseCase.GetChatByID(uint64(userID), uint64(chatID))
	assert.Equal(t, models.NewClientError(nil,
		http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//Test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{Members: []uint64{0}}, nil
		},
	}

	_, err = chatsUseCase.GetChatByID(uint64(userID), uint64(chatID))
	assert.Nil(t, err)

}

func TestChatsUseCaseImpl_GetChatsByUserID(t *testing.T) {
	userID := 0

	//Test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatsFunc: func(userID uint64) (chats []models.Chat, e error) {
			return []models.Chat{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.GetChatsByUserID(uint64(userID))

	assert.NotNil(t, err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatsFunc: func(userID uint64) (chats []models.Chat, e error) {
			return []models.Chat{models.Chat{Name: "mem", Members: []uint64{userID}}}, nil
		},
	}

	_, err = chatsUseCase.GetChatsByUserID(uint64(userID))

	assert.Nil(t, err)

	//test nil chats
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatsFunc: func(userID uint64) (chats []models.Chat, e error) {
			return []models.Chat{models.Chat{Name: "mem"}}, nil
		},
	}

	chats, err := chatsUseCase.GetChatsByUserID(uint64(userID))

	assert.Nil(t, err)
	assert.Equal(t, true, len(chats) == 0)
}

func TestChatsUseCaseImpl_DeleteChat(t *testing.T) {
	userID := 0
	chatID := 0

	//Internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.DeleteChat(uint64(userID), uint64(chatID))
	assert.NotNil(t, err)

	//Not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{Members: []uint64{1}}, nil
		},
	}

	err = chatsUseCase.DeleteChat(uint64(userID), uint64(chatID))
	assert.Equal(t, models.NewClientError(nil,
		http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//Test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChatByIDFunc: func(ID uint64) (chat models.Chat, e error) {
			return models.Chat{Members: []uint64{0}}, nil
		},
		RemoveChatFunc: func(chatID uint64) error {
			return nil
		},
	}

	err = chatsUseCase.DeleteChat(uint64(userID), uint64(chatID))
	assert.Nil(t, err)

}

func TestChatsUseCaseImpl_GetWorkspaceByID(t *testing.T) {
	chatID := 0
	userID := 0
	//Internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.GetWorkspaceByID(uint64(userID), uint64(chatID))
	assert.NotNil(t, err)

	//Not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{Members: []uint64{1}}, nil
		},
	}

	_, err = chatsUseCase.GetWorkspaceByID(uint64(userID), uint64(chatID))
	assert.Equal(t, models.NewClientError(nil,
		http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//Test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{Members: []uint64{0}}, nil
		},
	}

	_, err = chatsUseCase.GetWorkspaceByID(uint64(userID), uint64(chatID))
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_GetWorkspacesByUserID(t *testing.T) {
	userID := 0

	//Test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspacesFunc: func(userID uint64) (workspaces []models.Workspace, e error) {
			return []models.Workspace{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.GetWorkspacesByUserID(uint64(userID))

	assert.NotNil(t, err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspacesFunc: func(userID uint64) (workspaces []models.Workspace, e error) {
			return []models.Workspace{models.Workspace{Name: "mem", Members: []uint64{userID}}}, nil
		},
	}

	_, err = chatsUseCase.GetWorkspacesByUserID(uint64(userID))

	assert.Nil(t, err)

	//test nil chats
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspacesFunc: func(userID uint64) (workspaces []models.Workspace, e error) {
			return []models.Workspace{models.Workspace{Name: "mem"}}, nil
		},
	}

	chats, err := chatsUseCase.GetWorkspacesByUserID(uint64(userID))

	assert.Nil(t, err)
	assert.Equal(t, true, len(chats) == 0)
}

func TestChatsUseCaseImpl_EditWorkspace(t *testing.T) {
	userID := 0
	testWorkspace := &models.Workspace{ID: 0}

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.EditWorkspace(uint64(userID), testWorkspace)
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, nil
		},
	}

	err = chatsUseCase.EditWorkspace(uint64(userID), testWorkspace)
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{Admins: []uint64{uint64(userID)}}, nil
		},
		UpdateWorkspaceFunc: func(workspace *models.Workspace) error {
			return nil
		},
	}

	err = chatsUseCase.EditWorkspace(uint64(userID), testWorkspace)
	assert.Nil(t, err)

}

func TestChatsUseCaseImpl_LogoutFromWorkspace(t *testing.T) {
	userID := 0
	workspaceID := 0

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.LogoutFromWorkspace(uint64(userID), uint64(workspaceID))
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, nil
		},
	}

	err = chatsUseCase.LogoutFromWorkspace(uint64(userID), uint64(workspaceID))
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{Members: []uint64{uint64(userID)}}, nil
		},
		UpdateWorkspaceFunc: func(workspace *models.Workspace) error {
			return nil
		},
	}

	err = chatsUseCase.LogoutFromWorkspace(uint64(userID), uint64(workspaceID))
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_DeleteWorkspace(t *testing.T) {
	userID := 0
	workspaceID := 0

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.DeleteWorkspace(uint64(userID), uint64(workspaceID))
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{CreatorID: uint64(userID + 1)}, nil
		},
	}

	err = chatsUseCase.DeleteWorkspace(uint64(userID), uint64(workspaceID))
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{CreatorID: uint64(userID)}, nil
		},
		RemoveWorkspaceFunc: func(workspaceID uint64) error {
			return nil
		},
	}

	err = chatsUseCase.DeleteWorkspace(uint64(userID), uint64(workspaceID))
	assert.Nil(t, err)

}

func TestChatsUseCaseImpl_CreateChannel(t *testing.T) {
	userID := 0
	testChannel := &models.Channel{ID: 0}

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.CreateChannel(testChannel)
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{}, nil
		},
	}

	_, err = chatsUseCase.CreateChannel(testChannel)
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetWorkspaceByIDFunc: func(ID uint64) (workspace models.Workspace, e error) {
			return models.Workspace{Admins: []uint64{uint64(userID)}}, nil
		},
		PutChannelFunc: func(channel *models.Channel) (u uint64, e error) {
			return 0, nil
		},
	}

	_, err = chatsUseCase.CreateChannel(testChannel)
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_GetChannelByID(t *testing.T) {
	chatID := 0
	userID := 0
	//Internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, errors.New("Internal error")
		},
	}

	_, err := chatsUseCase.GetChannelByID(uint64(userID), uint64(chatID))
	assert.NotNil(t, err)

	//Not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{Members: []uint64{uint64(userID + 1)}}, nil
		},
	}

	_, err = chatsUseCase.GetChannelByID(uint64(userID), uint64(chatID))
	assert.Equal(t, models.NewClientError(nil,
		http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//Test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{Members: []uint64{uint64(userID)}}, nil
		},
	}

	_, err = chatsUseCase.GetChannelByID(uint64(userID), uint64(chatID))
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_EditChannel(t *testing.T) {
	userID := 0
	testChannel := &models.Channel{ID: 0}

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.EditChannel(uint64(userID), testChannel)
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, nil
		},
	}

	err = chatsUseCase.EditChannel(uint64(userID), testChannel)
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{Admins: []uint64{uint64(userID)}}, nil
		},
		UpdateChannelFunc: func(channel *models.Channel) error {
			return nil
		},
	}

	err = chatsUseCase.EditChannel(uint64(userID), testChannel)
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_DeleteChannel(t *testing.T) {
	userID := 0
	channelID := 0

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.DeleteChannel(uint64(userID), uint64(channelID))
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{CreatorID: uint64(userID + 1)}, nil
		},
	}

	err = chatsUseCase.DeleteChannel(uint64(userID), uint64(channelID))
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{CreatorID: uint64(userID)}, nil
		},
		RemoveChannelFunc: func(channelID uint64) error {
			return nil
		},
	}

	err = chatsUseCase.DeleteChannel(uint64(userID), uint64(channelID))
	assert.Nil(t, err)
}

func TestChatsUseCaseImpl_LogoutFromChannel(t *testing.T) {
	userID := 0
	channelID := 0

	//test internal error
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, errors.New("Internal error")
		},
	}

	err := chatsUseCase.LogoutFromChannel(uint64(userID), uint64(channelID))
	assert.NotNil(t, err)

	//test not permission
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{}, nil
		},
	}

	err = chatsUseCase.LogoutFromChannel(uint64(userID), uint64(channelID))
	assert.Equal(t,
		models.NewClientError(nil, http.StatusForbidden, "Not enough permissions for this request:("),
		err)

	//test success
	chatsUseCase.repository = &repository.ChatsRepositoryMock{
		GetChannelByIDFunc: func(ID uint64) (channel models.Channel, e error) {
			return models.Channel{Members: []uint64{uint64(userID)}}, nil
		},
		UpdateChannelFunc: func(channel *models.Channel) error {
			return nil
		},
	}

	err = chatsUseCase.LogoutFromChannel(uint64(userID), uint64(channelID))
	assert.Nil(t, err)
}
