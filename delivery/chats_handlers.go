package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

type ChatHandlers struct {
	Chats    useCase.ChatsUseCase
	Users    useCase.UsersUseCase
	Sessions repository.SessionRepository
}

func NewChatHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository, chats useCase.ChatsUseCase) ChatHandlers {
	return ChatHandlers{
		Chats:    chats,
		Users:    users,
		Sessions: sessions,
	}
}

func (handlers *ChatHandlers) sendError(err error, w http.ResponseWriter) {
	httpError, ok := err.(models.HTTPError)
	if !ok {
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := httpError.ResponseHeaders() // GetUserByEmail http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)

	_, err = w.Write(body)

	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (c *ChatHandlers) PostChat(w http.ResponseWriter, r *http.Request) {

	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newChatModel models.CreateChatModel
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChatModel)
	if err != nil {
		c.sendError(err, w)
		return
	}
	userTo, err := c.Users.GetUserByID(newChatModel.UserID)
	if err != nil {
		c.sendError(err, w)
		return
	}

	model := models.NewChatModel(userTo.Username, user.ID, userTo.ID)
	id, err := c.Chats.PutChat(model)
	if err != nil {
		c.sendError(err, w)
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
	chats, err := c.Chats.GetChatsByUserID(uint64(requestedID))
	if err != nil {

	}
	jsonChat, err := json.Marshal(chats)
	_, err = w.Write(jsonChat)
}

func (c *ChatHandlers) GetChatById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	chat, err := c.Chats.GetChatByID(user.ID, uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
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
		c.sendError(err, w)
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
		c.sendError(err, w)
		return
	}

	id, err := c.Chats.CreateChannel(&newChannelModel)
	if err != nil {
		c.sendError(err, w)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	channel, err := c.Chats.GetChannelByID(user.ID, uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
	}
	jsonChannel, err := json.Marshal(channel)
	_, err = w.Write(jsonChannel)
}

func (c *ChatHandlers) EditChannel(w http.ResponseWriter, r *http.Request) {

	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newChannel *models.Channel
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newChannel)
	if err != nil {
		c.sendError(err, w)
		return
	}

	err = c.Chats.EditChannel(user.ID, newChannel)
	if err != nil {
		c.sendError(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) RemoveChannel(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = c.Chats.DeleteChannel(user.ID, uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
	}
}

func (c *ChatHandlers) PostWorkspace(w http.ResponseWriter, r *http.Request) {
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var newWorkspace models.Workspace
	newWorkspace.Members = append(newWorkspace.Members, user.ID)
	newWorkspace.Admins = append(newWorkspace.Admins, user.ID)
	newWorkspace.CreatorID = user.ID
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newWorkspace)
	if err != nil {
		c.sendError(err, w)
		return
	}

	id, err := c.Chats.CreateWorkspace(&newWorkspace)
	if err != nil {
		c.sendError(err, w)
		return
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	if err != nil {
		c.sendError(err, w)
		return
	}
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers) GetWorkspaceById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := c.parseCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	workspace, err := c.Chats.GetWorkspaceByID(user.ID, uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
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
		c.sendError(err, w)
		return
	}

	err = c.Chats.EditWorkspace(user.ID, newWorkspace)
	if err != nil {
		c.sendError(err, w)
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
		c.sendError(err, w)
	}
}

func (c ChatHandlers) parseCookie(r *http.Request) (models.User, error) {
	cookie, _ := r.Cookie("session_id")
	id, err := c.Sessions.GetID(cookie.Value)
	user, err := c.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}
