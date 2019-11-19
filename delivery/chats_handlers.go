package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ChatHandlers struct {
	Chats    useCase.ChatsUseCase
	Users    useCase.UsersUseCase
	Sessions repository.SessionRepository
	utils    utils.HandlersUtils
}

func NewChatHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository,
	chats useCase.ChatsUseCase, utils utils.HandlersUtils) ChatHandlers {
	return ChatHandlers{
		Chats:    chats,
		Users:    users,
		Sessions: sessions,
		utils:    utils,
	}
}

func (c *ChatHandlers) PostChat(w http.ResponseWriter, r *http.Request) {

	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	var newChatModel models.CreateChatModel
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChatModel)
	if err != nil {
		c.utils.HandleError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
		return
	}
	userTo, err := c.Users.GetUserByID(newChatModel.UserID)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}

	model := models.NewChatModel(userTo.Username, user.ID, userTo.ID)
	id, err := c.Chats.PutChat(model)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) GetChatsByUser(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	cookie, _ := r.Cookie("session_id")
	cookieID, err := c.Sessions.GetID(cookie.Value)
	if err != nil {
		c.utils.HandleError(
			models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:("),
			w, r)
		return
	}
	if cookieID != uint64(requestedID) {
		c.utils.HandleError(
			models.NewClientError(err, http.StatusUnauthorized, fmt.Sprintf("Actual id: %d, Requested id: %d", cookieID, requestedID)),
			w, r)
		return
	}
	chats, err := c.Chats.GetChatsByUserID(uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	workspaces, err := c.Chats.GetWorkspacesByUserID(uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	responseChats := models.ResponseChatsArray{Chats: chats, Workspaces: workspaces}
	jsonChat, err := json.Marshal(responseChats)
	_, err = w.Write(jsonChat)
}

func (c *ChatHandlers) GetChatById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	chat, err := c.Chats.GetChatByID(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	jsonChat, err := json.Marshal(chat)
	_, err = w.Write(jsonChat)
}

func (c *ChatHandlers) RemoveChat(w http.ResponseWriter, r *http.Request) {

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	//TODO:Check error
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = c.Chats.DeleteChat(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
}

func (c *ChatHandlers) PostChannel(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newChannelModel models.Channel
	newChannelModel.Members = append(newChannelModel.Members, user.ID)
	newChannelModel.Admins = append(newChannelModel.Admins, user.ID)
	newChannelModel.CreatorID = user.ID
	newChannelModel.WorkspaceID = uint64(requestedID)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChannelModel)
	if err != nil {
		c.utils.HandleError(models.NewClientError(err, http.StatusBadRequest,
			"Bad request: malformed data:("), w, r)
		return
	}

	id, err := c.Chats.CreateChannel(&newChannelModel)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) GetChannelById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	channel, err := c.Chats.GetChannelByID(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	jsonChannel, err := json.Marshal(channel)
	_, err = w.Write(jsonChannel)
}

func (c *ChatHandlers) EditChannel(w http.ResponseWriter, r *http.Request) {
	requestedID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		c.utils.HandleError(models.NewClientError(nil, http.StatusBadRequest, ""), w, r)
	}
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	var newChannel *models.Channel
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChannel)
	newChannel.ID = uint64(requestedID)
	if err != nil {
		c.utils.HandleError(models.NewClientError(err, http.StatusBadRequest,
			"Bad request: malformed data:("), w, r)
		return
	}

	err = c.Chats.EditChannel(user.ID, newChannel)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) RemoveChannel(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	err = c.Chats.DeleteChannel(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
}

func (c *ChatHandlers) PostWorkspace(w http.ResponseWriter, r *http.Request) {
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	var newWorkspace models.Workspace
	newWorkspace.Members = append(newWorkspace.Members, user.ID)
	newWorkspace.Admins = append(newWorkspace.Admins, user.ID)
	newWorkspace.CreatorID = user.ID
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newWorkspace)
	if err != nil {
		c.utils.HandleError(models.NewClientError(err, http.StatusBadRequest,
			"Bad request: malformed data:("), w, r)
		return
	}

	id, err := c.Chats.CreateWorkspace(&newWorkspace)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	newWorkspace.ID = id
	jsonResponse, err := json.Marshal(newWorkspace)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) GetWorkspaceById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	workspace, err := c.Chats.GetWorkspaceByID(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	jsonWorkspace, err := json.Marshal(workspace)
	_, err = w.Write(jsonWorkspace)
}

func (c *ChatHandlers) EditWorkspace(w http.ResponseWriter, r *http.Request) {
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newWorkspace *models.Workspace
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newWorkspace)
	if err != nil {
		c.utils.HandleError(models.NewClientError(err, http.StatusBadRequest,
			"Bad request: malformed data:("), w, r)
		return
	}

	err = c.Chats.EditWorkspace(user.ID, newWorkspace)
	if err != nil {
		c.utils.HandleError(err, w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) RemoveWorkspace(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = c.Chats.DeleteWorkspace(user.ID, uint64(requestedID))
	if err != nil {
		c.utils.HandleError(err, w, r)
	}
}

func (c ChatHandlers) parseCookie(r *http.Request) (models.User, error) {
	cookie, _ := r.Cookie("session_id")
	id, err := c.Sessions.GetID(cookie.Value)
	if err != nil {
		return models.User{}, models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
	user, err := c.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, err
	}
}
