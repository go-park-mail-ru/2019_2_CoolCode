package delivery

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	"io/ioutil"
	"net/http"
	"testing"
)

var messageApi MessageHandlersImpl

type MessageTestCase struct {
	name         string
	Body         string
	SessionID    string
	Headers      map[string]string
	Method       string
	URL          string
	Response     string
	StatusCode   int
	Handler      http.HandlerFunc
	Messages     useCase.MessagesUseCase
	Notification useCase.NotificationUseCase
	Users        useCase.UsersUseCase
	Sessions     repository.SessionRepository
	utils        utils.HandlersUtils
}

func runTableMessageAPITests(t *testing.T, cases []*MessageTestCase) {
	for _, c := range cases {
		runMessageAPITest(t, c)
	}
}

func runMessageAPITest(t *testing.T, testCase *MessageTestCase) {
	t.Run(testCase.name, func(t *testing.T) {
		if testCase.Messages != nil {
			messageApi.Messages = testCase.Messages
		}
		if testCase.Notification != nil {
			messageApi.Notifications = testCase.Notification
		}
		if testCase.Sessions != nil {
			messageApi.Sessions = testCase.Sessions
		}
		if testCase.Users != nil {
			messageApi.Users = testCase.Users
		}
		apitest.New().
			Handler(testCase.Handler).
			Method(testCase.Method).
			Headers(testCase.Headers).
			Cookie("session_id", "test").
			URL(testCase.URL).
			Body(testCase.Body).
			Expect(t).
			Status(testCase.StatusCode).End()
	})
}

func init() {
	emptyLogger := logrus.New()
	emptyLogger.Out = ioutil.Discard
	messageApi.utils = utils.NewHandlersUtils(emptyLogger)
}

func TestMessageHandlersImpl_SendMessage(t *testing.T) {
	tests := []*MessageTestCase{
		{
			name: "WrongCookieTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Internal error")
				},
			},
			StatusCode: 401,
			Handler:    http.HandlerFunc(messageApi.SendMessage),
		},
		{
			name: "InternalErrorTestGetUserByID",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, errors.New("Internal error")
				},
			},
			StatusCode: 401,
			Handler:    http.HandlerFunc(messageApi.SendMessage),
		},
		{
			name: "InvalidJSONTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Body:       "Bad JSON",
			StatusCode: 400,
			Handler:    http.HandlerFunc(messageApi.SendMessage),
		},
		{
			name: "InternalSaveMessageErrorTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Messages: &useCase.MessagesUseCaseMock{
				SaveMessageFunc: func(message *models.Message) (u uint64, e error) {
					return 0, errors.New("Internal error")
				},
			},
			Body:       `{"text":"mem"}`,
			StatusCode: 500,
			Handler:    http.HandlerFunc(messageApi.SendMessage),
		},
		{
			name: "SendMessageErrorTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Messages: &useCase.MessagesUseCaseMock{
				SaveMessageFunc: func(message *models.Message) (u uint64, e error) {
					return 0, nil
				},
			},
			Notification: &useCase.NotificationUseCaseMock{
				SendMessageFunc: func(chatID uint64, message []byte) error {
					return errors.New("Internal error")
				},
			},
			Body:       `{"text":"mem"}`,
			StatusCode: 200,
			Handler:    http.HandlerFunc(messageApi.SendMessage),
		},
	}
	runTableMessageAPITests(t, tests)
}

func TestMessageHandlersImpl_EditMessage(t *testing.T) {
	tests := []*MessageTestCase{
		{
			name: "WrongCookieTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Internal error")
				},
			},
			StatusCode: 401,
			Handler:    http.HandlerFunc(messageApi.EditMessage),
		},
		{
			name: "InvalidJSONTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Body:       "Bad JSON",
			StatusCode: 400,
			Handler:    http.HandlerFunc(messageApi.EditMessage),
		},
		{
			name: "InternalEditMessageErrorTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Messages: &useCase.MessagesUseCaseMock{
				EditMessageFunc: func(message *models.Message, userID uint64) error {
					return errors.New("Internal error")
				},
			},
			Body:       `{"text":"mem"}`,
			StatusCode: 500,
			Handler:    http.HandlerFunc(messageApi.EditMessage),
		},
	}

	runTableMessageAPITests(t, tests)

}

func TestMessageHandlersImpl_GetMessagesByChatID(t *testing.T) {
	tests := []*MessageTestCase{
		{
			name: "WrongCookieTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Internal error")
				},
			},
			StatusCode: 401,
			Handler:    http.HandlerFunc(messageApi.GetMessagesByChatID),
		},
		{
			name: "InternalGetMessageErrorTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Messages: &useCase.MessagesUseCaseMock{
				GetChatMessagesFunc: func(chatID uint64, userID uint64) (messages models.Messages, e error) {
					return models.Messages{}, errors.New("Internal error")
				},
			},
			Body:       `{"text":"mem"}`,
			StatusCode: 500,
			Handler:    http.HandlerFunc(messageApi.GetMessagesByChatID),
		},
		{
			name: "SuccessTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Messages: &useCase.MessagesUseCaseMock{
				GetChatMessagesFunc: func(chatID uint64, userID uint64) (messages models.Messages, e error) {
					return models.Messages{}, nil
				},
			},
			Body:       `{"text":"mem"}`,
			StatusCode: 200,
			Handler:    http.HandlerFunc(messageApi.GetMessagesByChatID),
		},
	}

	runTableMessageAPITests(t, tests)
}
