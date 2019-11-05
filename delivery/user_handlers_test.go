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
	"os"
	"testing"
)

var usersApi UserHandlers

type UserTestCase struct {
	name       string
	Body       string
	SessionID  string
	Headers    map[string]string
	Method     string
	URL        string
	Response   string
	StatusCode int
	Handler    http.HandlerFunc
	Users      useCase.UsersUseCase
	Photos     repository.PhotoRepository
	Sessions   repository.SessionRepository
	utils      utils.HandlersUtils
}

func runTableAPITests(t *testing.T, cases []*UserTestCase) {
	for _, c := range cases {
		runAPITest(t, c)
	}
}

func runAPITest(t *testing.T, testCase *UserTestCase) {
	t.Run(testCase.name, func(t *testing.T) {
		if testCase.Users != nil {
			usersApi.Users = testCase.Users
		}
		if testCase.Sessions != nil {
			usersApi.Sessions = testCase.Sessions
		}
		if testCase.Photos != nil {
			usersApi.Photos = testCase.Photos
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
	usersApi.utils = utils.NewHandlersUtils(emptyLogger)
}

func TestUserHandlers_SignUp(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "TestBadJson",
			Body:       "asfasfasfsadf",
			Method:     "Post",
			StatusCode: 400,
			Handler:    http.HandlerFunc(usersApi.SignUp),
			URL:        "/users",
		},
		{
			name:       "TestBadJson",
			Body:       "",
			Handler:    http.HandlerFunc(usersApi.SignUp),
			Method:     "Post",
			StatusCode: 400,
			URL:        "/users",
		},
		{
			name:       "TestUserExists",
			Body:       `{"id": 1}`,
			Method:     "Post",
			Handler:    http.HandlerFunc(usersApi.SignUp),
			StatusCode: 400,
			URL:        "/users",
			Users: &useCase.UsersUseCaseMock{
				SignUpFunc: func(user *models.User) error {
					return models.NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
				},
			},
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_Login(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "TestBadJson",
			Body:       "mem",
			Method:     "POST",
			StatusCode: 400,
			Handler:    http.HandlerFunc(usersApi.Login),
			URL:        "login",
		},
		{
			name:       "TestWrongPassword",
			Body:       `{"id":0}`,
			Method:     "POST",
			StatusCode: 400,
			Handler:    http.HandlerFunc(usersApi.Login),
			Users: &useCase.UsersUseCaseMock{
				LoginFunc: func(user models.User) (user2 models.User, e error) {
					return models.User{}, models.NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
				},
			},
			URL: "login",
		},
		{
			name:       "TestCookieError",
			Body:       `{"id":0}`,
			Method:     "POST",
			StatusCode: 500,
			Handler:    http.HandlerFunc(usersApi.Login),
			Users: &useCase.UsersUseCaseMock{
				LoginFunc: func(user models.User) (user2 models.User, e error) {
					return models.User{}, nil
				},
			},
			Sessions: &repository.SessionRepositoryMock{
				PutFunc: func(session string, id uint64) error {
					return errors.New("Internal error")
				},
			},
			URL: "login",
		},
		{
			name:       "TestSuccess",
			Body:       `{"id":0}`,
			Method:     "POST",
			StatusCode: 200,
			Handler:    http.HandlerFunc(usersApi.Login),
			Users: &useCase.UsersUseCaseMock{
				LoginFunc: func(user models.User) (user2 models.User, e error) {
					return models.User{}, nil
				},
			},
			Sessions: &repository.SessionRepositoryMock{
				PutFunc: func(session string, id uint64) error {
					return nil
				},
			},
			URL: "login",
		},
	}

	runTableAPITests(t, tests)
}

func TestUserHandlers_SavePhoto(t *testing.T) {
	tests := []*UserTestCase{
		{
			name: "WrongCookieTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, models.NewClientError(nil, http.StatusUnauthorized, "Bad request : not valid cookie:(")
				},
			},
			Handler:    http.HandlerFunc(usersApi.SavePhoto),
			StatusCode: 401,
		},
		{
			name: "WrongDataTypeTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Handler:    http.HandlerFunc(usersApi.SavePhoto),
			StatusCode: 500,
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_GetPhoto(t *testing.T) {
	tests := []*UserTestCase{
		{
			name: "InternalErrorTest",
			Photos: &repository.PhotoRepositoryMock{
				GetPhotoFunc: func(id int) (file *os.File, e error) {
					return &os.File{}, errors.New("Internal error")
				},
			},
			URL:        "photos/1",
			Handler:    http.HandlerFunc(usersApi.GetPhoto),
			StatusCode: 500,
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_GetUser(t *testing.T) {
	tests := []*UserTestCase{
		{
			name: "WrongCookieTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, models.NewClientError(nil, http.StatusUnauthorized, "Bad request : not valid cookie:(")
				},
			},
			URL:        "users/1",
			Handler:    http.HandlerFunc(usersApi.GetUser),
			StatusCode: 401,
		},
		{
			name: "InternalError",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 1, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					if user.ID == 0 {
						return models.User{}, errors.New("Internal error")
					} else {
						return models.User{}, nil
					}
				},
			},
			URL:        "users/1",
			Handler:    http.HandlerFunc(usersApi.GetUser),
			StatusCode: 500,
		},
		{
			name: "SuccessTest",
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {

					return models.User{}, nil
				},
			},
			URL:        "users/1",
			Handler:    http.HandlerFunc(usersApi.GetUser),
			StatusCode: 200,
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_GetUserBySession(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "WrongCookieTest",
			Handler:    http.HandlerFunc(usersApi.GetUserBySession),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Wrong cookie")
				},
			},
			URL:    "/users",
			Method: "GET",
		},
		{
			name:       "SuccessTest",
			Handler:    http.HandlerFunc(usersApi.GetUserBySession),
			StatusCode: 200,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			URL:    "/users",
			Method: "GET",
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_EditProfile(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "WrongCookieTest",
			Handler:    http.HandlerFunc(usersApi.EditProfile),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Wrong cookie")
				},
			},
		},
		{
			name:       "WrongUserIDTest",
			Handler:    http.HandlerFunc(usersApi.EditProfile),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{ID: 1}, nil
				},
			},
		},
		{
			name:       "WrongUserIDTest",
			Handler:    http.HandlerFunc(usersApi.EditProfile),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Body: `{"id":1}`,
		},
		{
			name:       "BadJsonTest",
			Handler:    http.HandlerFunc(usersApi.EditProfile),
			StatusCode: 400,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
			Body: `BadJson`,
		},
		{
			name:       "InternalErrorTest",
			Handler:    http.HandlerFunc(usersApi.EditProfile),
			StatusCode: 500,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
				ChangeUserFunc: func(user *models.User) error {
					return errors.New("Internal error")
				},
			},

			Body: `{"id":0}`,
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_Logout(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "WrongCookieTest",
			Handler:    http.HandlerFunc(usersApi.Logout),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				RemoveFunc: func(session string) error {
					return errors.New("InternalError")
				},
			},
		},
		{
			name:       "SuccessTest",
			Handler:    http.HandlerFunc(usersApi.Logout),
			StatusCode: 200,
			Sessions: &repository.SessionRepositoryMock{
				RemoveFunc: func(session string) error {
					return nil
				},
			},
		},
	}
	runTableAPITests(t, tests)
}

func TestUserHandlers_FindUsers(t *testing.T) {
	tests := []*UserTestCase{
		{
			name:       "WrongCookieTest",
			Handler:    http.HandlerFunc(usersApi.FindUsers),
			StatusCode: 401,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, errors.New("Wrong cookie")
				},
			},
		},
		{
			name:       "InternalErrorTest",
			Handler:    http.HandlerFunc(usersApi.FindUsers),
			StatusCode: 500,
			Sessions: &repository.SessionRepositoryMock{
				GetIDFunc: func(session string) (u uint64, e error) {
					return 0, nil
				},
			},
			Users: &useCase.UsersUseCaseMock{
				FindUsersFunc: func(name string) (users models.Users, e error) {
					return models.Users{}, errors.New("Internal Error")
				},
				GetUserByIDFunc: func(id uint64) (user models.User, e error) {
					return models.User{}, nil
				},
			},
		},
	}
	runTableAPITests(t, tests)
}
