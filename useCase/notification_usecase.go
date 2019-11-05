package useCase

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
)

//go:generate moq -out notifications_ucase_mock.go . NotificationUseCase
type NotificationUseCase interface {
	OpenConn(ID uint64) (*models.WebSocketHub, error)
	SendMessage(chatID uint64, message []byte) error
}

type NotificationUseCaseImpl struct {
	notificationRepository repository.NotificationRepository
}

func NewNotificationUseCase() NotificationUseCase {
	return &NotificationUseCaseImpl{notificationRepository: repository.NewArrayRepo()}
}

func (u *NotificationUseCaseImpl) OpenConn(ID uint64) (*models.WebSocketHub, error) {
	return u.notificationRepository.GetNotificationHub(ID), nil
}

func (u *NotificationUseCaseImpl) SendMessage(chatID uint64, message []byte) error {
	hub := u.notificationRepository.GetNotificationHub(chatID)
	hub.BroadcastChan <- message
	return nil
}
