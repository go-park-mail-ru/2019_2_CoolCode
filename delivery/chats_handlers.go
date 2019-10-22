package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/mux"
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
	cookie, _ := r.Cookie("session_id")

	user, err := c.parseCookie(cookie)
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
	err = c.Chats.PutChat(model)
	if err != nil {
		c.sendError(err, w)
		return
	}
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
	chat, err := c.Chats.GetChatByID(uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
	}
	jsonChat, err := json.Marshal(chat)
	_, err = w.Write(jsonChat)
}

func (c *ChatHandlers) GetChannelById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	channel, err := c.Chats.GetChannelByID(uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
	}
	jsonChannel, err := json.Marshal(channel)
	_, err = w.Write(jsonChannel)
}

func (c *ChatHandlers) GetWorkspaceById(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	workspace, err := c.Chats.GetWorkspaceByID(uint64(requestedID))
	if err != nil {
		c.sendError(err, w)
	}
	jsonWorkspace, err := json.Marshal(workspace)
	_, err = w.Write(jsonWorkspace)
}

func (c ChatHandlers) parseCookie(cookie *http.Cookie) (models.User, error) {
	id, err := c.Sessions.GetID(cookie.Value)
	user, err := c.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}
