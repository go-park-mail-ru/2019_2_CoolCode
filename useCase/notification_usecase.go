package useCase

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
)

type NotificationUseCase struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationUseCase() NotificationUseCase {
	return NotificationUseCase{notificationRepository: repository.NewArrayRepo()}
}

func (u *NotificationUseCase) OpenConn(ID uint64) (*models.WebSocketHub, error) {
	return u.notificationRepository.GetNotificationHub(ID), nil
}

func (u *NotificationUseCase) SendMessage(chatID uint64, message []byte) error {
	hub := u.notificationRepository.GetNotificationHub(chatID)
	hub.BroadcastChan <- message
	return nil
}
