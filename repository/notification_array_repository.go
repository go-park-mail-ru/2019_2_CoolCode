package repository

import "github.com/go-park-mail-ru/2019_2_CoolCode/models"

type NotificationArrayRepository struct {
	Hubs map[uint64]*models.WebSocketHub
}

func (n *NotificationArrayRepository) GetNotificationHub(chatID uint64) *models.WebSocketHub {
	if hub, ok := n.Hubs[chatID]; ok {
		return hub
	}
	n.Hubs[chatID] = models.NewHub()
	return n.Hubs[chatID]
}

func NewArrayRepo() NotificationRepository {
	return &NotificationArrayRepository{Hubs: make(map[uint64]*models.WebSocketHub, 0)}
}
