package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type MessageHandlers interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetMessagesByChatID(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
	EditMessage(w http.ResponseWriter, r *http.Request)
}

type MessageHandlersImpl struct {
	useCase       useCase.MessagesUseCase
	Users         useCase.UsersUseCase
	Sessions      repository.SessionRepository
	Notifications useCase.NotificationUseCase
}

func NewMessageHandlers(useCase useCase.MessagesUseCase, users useCase.UsersUseCase, sessions repository.SessionRepository, notificationUseCase useCase.NotificationUseCase) MessageHandlers {
	return &MessageHandlersImpl{
		useCase:       useCase,
		Users:         users,
		Sessions:      sessions,
		Notifications: notificationUseCase,
	}
}

func (m *MessageHandlersImpl) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		utils.SendError(err, w)
		return
	}
	message, err := parseMessage(r)
	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	message.AuthorID = user.ID
	message.ChatID = uint64(chatID)
	id, err := m.useCase.SaveMessage(message)
	if err != nil {
		utils.SendError(err, w)
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("An error occurred %v", err)
	}

	//send to websocket
	message.ID = id
	websocketMessage := models.WebsocketMessage{
		WebsocketEventType: 1,
		Body:               *message,
	}
	websocketJson, err := json.Marshal(websocketMessage)
	if err != nil {
		log.Printf("An error occurred %v", err)
	}
	m.Notifications.SendMessage(message.ChatID, websocketJson)

}

func (m *MessageHandlersImpl) EditMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		utils.SendError(err, w)
		return
	}
	message, err := parseMessage(r)

	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	message.ID = uint64(messageID)

	err = m.useCase.EditMessage(message, user.ID)
	if err != nil {
		utils.SendError(err, w)
	}
}

func (m *MessageHandlersImpl) GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		utils.SendError(err, w)
		return
	}
	messages, err := m.useCase.GetChatMessages(uint64(chatID), user.ID)
	if err != nil {
		utils.SendError(err, w)
	}
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		utils.SendError(err, w)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("An error occurred %v", err)
	}
}

func (m *MessageHandlersImpl) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.SendError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		utils.SendError(err, w)
		return
	}

	hide, ok := r.URL.Query()["forAuthor"]
	if !ok || len(hide[0]) < 1 {
		err = m.useCase.DeleteMessage(uint64(messageID), user.ID)
	} else {
		err = m.useCase.HideMessageForAuthor(uint64(messageID), user.ID)
	}

	if err != nil {
		utils.SendError(err, w)
	}
}

func (m *MessageHandlersImpl) parseCookie(r *http.Request) (models.User, error) {
	cookie, _ := r.Cookie("session_id")
	id, err := m.Sessions.GetID(cookie.Value)
	user, err := m.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}

func parseMessage(r *http.Request) (*models.Message, error) {
	var message models.Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	return &message, err
}
