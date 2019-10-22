package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ChatHandlers struct {
	Chats    useCase.ChatsUseCase
	Users    useCase.UsersUseCase
	Sessions repository.SessionRepository
}

func NewChatHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository) *ChatHandlers {
	return &ChatHandlers{
		Chats:    useCase.NewChatsUseCase(repository.NewChatArrayRepository()),
		Users:    users,
		Sessions: sessions,
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
		//TODO: send error
	}
	userTo, err := c.Users.GetUserByID(newChatModel.UserID)
	if err != nil {
		//TODO: send error
	}

	model := models.NewChatModel(userTo.Username, user.ID, userTo.ID)
	err = c.Chats.PutChat(model)
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

func (c *ChatHandlers) parseCookie(cookie *http.Cookie) (models.User, error) {
	id, err := c.Sessions.GetID(cookie.Value)
	user, err := c.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}
