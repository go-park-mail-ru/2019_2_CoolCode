package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type MessageHandlers interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetMessagesByChatID(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
	EditMessage(w http.ResponseWriter, r *http.Request)
	FindMessages(w http.ResponseWriter, r *http.Request)
}

type MessageHandlersImpl struct {
	Messages      useCase.MessagesUseCase
	Users         useCase.UsersUseCase
	Sessions      repository.SessionRepository
	Notifications useCase.NotificationUseCase
	utils         utils.HandlersUtils
}

func NewMessageHandlers(useCase useCase.MessagesUseCase, users useCase.UsersUseCase,
	sessions repository.SessionRepository, notificationUseCase useCase.NotificationUseCase, handlersUtils utils.HandlersUtils) MessageHandlers {
	return &MessageHandlersImpl{
		Messages:      useCase,
		Users:         users,
		Sessions:      sessions,
		Notifications: notificationUseCase,
		utils:         handlersUtils,
	}
}

func (m *MessageHandlersImpl) SendMessage(w http.ResponseWriter, r *http.Request) {
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])

	var id uint64
	if err != nil {
		m.utils.LogError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), r)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	message, err := parseMessage(r)
	if err != nil {
		m.utils.HandleError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
		return
	}
	message.AuthorID = user.ID
	message.ChatID = uint64(chatID)
	if isChannel(r) {
		id, err = m.Messages.SaveChannelMessage(message)
	} else {
		id, err = m.Messages.SaveChatMessage(message)
	}
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(map[string]uint64{
		"id": id,
	})
	_, err = w.Write(jsonResponse)
	if err != nil {
		m.utils.LogError(err, r)
	}

	//send to websocket
	message.ID = id
	websocketMessage := models.WebsocketMessage{
		WebsocketEventType: 1,
		Body:               *message,
	}
	websocketJson, err := json.Marshal(websocketMessage)
	if err != nil {
		m.utils.LogError(err, r)
	}
	err = m.Notifications.SendMessage(message.ChatID, websocketJson)
	if err != nil {
		m.utils.LogError(err, r)
	}

}

func (m *MessageHandlersImpl) EditMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		m.utils.LogError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), r)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	message, err := parseMessage(r)

	if err != nil {
		m.utils.HandleError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
		return
	}
	message.ID = uint64(messageID)
	dbMessage, err := m.Messages.GetMessageByID(uint64(messageID))
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	message.ChatID = dbMessage.ChatID

	err = m.Messages.EditMessage(message, user.ID)
	if err != nil {
		m.utils.HandleError(err, w, r)
	}

	//send to websocket

	websocketMessage := models.WebsocketMessage{
		WebsocketEventType: 3,
		Body:               *message,
	}
	websocketJson, err := json.Marshal(websocketMessage)
	if err != nil {
		m.utils.LogError(err, r)
	}
	err = m.Notifications.SendMessage(message.ChatID, websocketJson)
	if err != nil {
		m.utils.LogError(err, r)
	}
}

func (m *MessageHandlersImpl) GetMessagesByChatID(w http.ResponseWriter, r *http.Request) {
	var messages models.Messages
	chatID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		m.utils.LogError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), r)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	if isChannel(r) {
		messages, err = m.Messages.GetChannelMessages(uint64(chatID), user.ID)
	} else {
		messages, err = m.Messages.GetChatMessages(uint64(chatID), user.ID)
	}
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		m.utils.HandleError(err, w, r)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		m.utils.LogError(err, r)
	}
}

func (m *MessageHandlersImpl) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		m.utils.HandleError(models.NewClientError(err, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}

	message, err := m.Messages.GetMessageByID(uint64(messageID))
	if err != nil {
		m.utils.LogError(err, r)
	}

	hide, ok := r.URL.Query()["forAuthor"]
	if !ok || len(hide[0]) < 1 {
		err = m.Messages.DeleteMessage(uint64(messageID), user.ID)
	} else {
		err = m.Messages.HideMessageForAuthor(uint64(messageID), user.ID)
	}

	if err != nil {
		m.utils.HandleError(err, w, r)
	}

	websocketMessage := models.WebsocketMessage{
		WebsocketEventType: 2,
		Body:               *message,
	}
	websocketJson, err := json.Marshal(websocketMessage)
	if err != nil {
		m.utils.LogError(err, r)
	}
	err = m.Notifications.SendMessage(message.ChatID, websocketJson)
	if err != nil {
		m.utils.LogError(err, r)
	}
}

func (m *MessageHandlersImpl) parseCookie(r *http.Request) (models.User, error) {
	cookie, _ := r.Cookie("session_id")
	id, err := m.Sessions.GetID(cookie.Value)
	if err != nil {
		return models.User{}, models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
	user, err := m.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}

func (m *MessageHandlersImpl) FindMessages(w http.ResponseWriter, r *http.Request) {
	findString, ok := mux.Vars(r)["text"]
	if !ok {
		m.utils.HandleError(models.NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data:("), w, r)
	}
	user, err := m.parseCookie(r)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}
	messages, err := m.Messages.FindMessages(findString, user.ID)
	if err != nil {
		m.utils.HandleError(err, w, r)
		return
	}

	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		m.utils.HandleError(err, w, r)
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		m.utils.LogError(err, r)
	}

}

func parseMessage(r *http.Request) (*models.Message, error) {
	var message models.Message
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	return &message, err
}

func isChannel(r *http.Request) bool {
	return strings.Contains(r.URL.String(), "channels")
}
