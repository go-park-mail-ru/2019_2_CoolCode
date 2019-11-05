package delivery

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type NotificationHandlers struct {
	notificationUseCase useCase.NotificationUseCase
	chatsUseCase        useCase.ChatsUseCase
	Users               useCase.UsersUseCase
	Sessions            repository.SessionRepository
	utils               utils.HandlersUtils
}

func NewNotificationHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository,
	chats useCase.ChatsUseCase, notifications useCase.NotificationUseCase, utils utils.HandlersUtils) NotificationHandlers {
	return NotificationHandlers{
		notificationUseCase: notifications,
		chatsUseCase:        chats,
		Users:               users,
		Sessions:            sessions,
		utils:               utils,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *NotificationHandlers) HandleNewWSConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.utils.HandleError(models.NewServerError(err, http.StatusBadRequest, "Can not upgrade connection"), w, r)
		return
	}
	sessionID, _ := r.Cookie("session_id")

	requestedID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.utils.LogError(err, r)
	}

	userID, err := h.parseCookie(sessionID)
	if err != nil {
		h.utils.HandleError(err, w, r)
		return
	}
	//Проверяем доступ к чату
	ok, err := h.chatsUseCase.CheckChatPermission(userID, uint64(requestedID))
	if err != nil {
		h.utils.HandleError(err, w, r)
		return
	}
	if !ok {
		h.utils.HandleError(models.NewClientError(nil, http.StatusForbidden, "Not permission to chat:("),
			w, r)
		return
	}
	//Достаем Handler с помощью Messages
	hub, err := h.notificationUseCase.OpenConn(uint64(requestedID))
	go hub.Run()
	//Запускаем event loop
	hub.AddClientChan <- ws

	for {
		var m []byte

		_, m, err := ws.ReadMessage()

		if err != nil {
			hub.RemoveClient(ws)
			return
		}
		hub.BroadcastChan <- m
	}

}

func (h NotificationHandlers) parseCookie(cookie *http.Cookie) (uint64, error) {
	ID, err := h.Sessions.GetID(cookie.Value)
	if err == nil {
		return ID, nil
	} else {
		return ID, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
}
