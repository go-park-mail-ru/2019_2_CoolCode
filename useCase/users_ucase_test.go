package useCase

import (
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var useCase = usersUseCase{
	repository: &repository.UserRepoMock{},
}

func TestUsersUseCase_Login(t *testing.T) {
	testUser := models.User{
		Username: "",
		Email:    "mem@mem.ru",
		Name:     "",
		Password: "mocktestpassword",
		Status:   "",
		Phone:    "",
	}

	//user not exists
	useCase.repository = &repository.UserRepoMock{
		GetUserByEmailFunc: func(email string) (user models.User, e error) {
			return models.User{}, errors.New("Not contains")
		},
	}
	_, err := useCase.Login(testUser)
	assert.NotNil(t, err)

	//wrong password
	useCase.repository = &repository.UserRepoMock{
		GetUserByEmailFunc: func(email string) (user models.User, e error) {
			return models.User{}, nil
		},
	}

	_, err = useCase.Login(testUser)
	assert.NotNil(t, err)

	//success
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.MinCost)
	if err != nil {
		t.Error(err)
	}
	useCase.repository = &repository.UserRepoMock{
		GetUserByEmailFunc: func(email string) (user models.User, e error) {
			return models.User{Password: string(hashedPassword)}, nil
		},
	}

	_, err = useCase.Login(testUser)
	assert.Nil(t, err)
}

func TestUsersUseCase_SignUp(t *testing.T) {
	testUser := &models.User{
		Email:    "mem@mem.ru",
		Password: "1",
		Username: "mem",
	}

	//test contains user
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return true
		},
	}

	err := useCase.SignUp(testUser)
	assert.NotNil(t, err)

	//test internal error
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return false
		},
		PutUserFunc: func(newUser *models.User) (u uint64, e error) {
			return 0, errors.New("Internal error")
		},
	}

	err = useCase.SignUp(testUser)
	assert.NotNil(t, err)

	//test success
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return false
		},
		PutUserFunc: func(newUser *models.User) (u uint64, e error) {
			return 0, nil
		},
	}

	err = useCase.SignUp(testUser)
	assert.Nil(t, err)

}

func TestUsersUseCase_GetUserByID(t *testing.T) {
	testID := 0

	//test err
	useCase.repository = &repository.UserRepoMock{
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, errors.New("repository err")
		},
	}

	_, err := useCase.GetUserByID(uint64(testID))
	assert.NotNil(t, err)

	//test not valid user
	useCase.repository = &repository.UserRepoMock{
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, nil
		},
	}

	_, err = useCase.GetUserByID(uint64(testID))
	assert.NotNil(t, err)

	//test success
	useCase.repository = &repository.UserRepoMock{
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{Email: "mem@mem.ru"}, nil
		},
	}

	_, err = useCase.GetUserByID(uint64(testID))
	assert.Nil(t, err)
}

func TestUsersUseCase_GetUserByEmail(t *testing.T) {
	testEmail := "mem@mem.ru"

	useCase.repository = &repository.UserRepoMock{
		GetUserByEmailFunc: func(email string) (user models.User, e error) {
			return models.User{}, nil
		},
	}

	_, err := useCase.GetUserByEmail(testEmail)
	assert.Nil(t, err)
}

func TestUsersUseCase_ChangeUser(t *testing.T) {
	testUser := &models.User{
		Email:    "mem@mem.ru",
		Username: "mem",
	}

	//test contains user
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return false
		},
	}

	err := useCase.ChangeUser(testUser)
	assert.NotNil(t, err)

	//test internal error GetByID
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return true
		},
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, errors.New("Internal error")
		},
	}

	err = useCase.ChangeUser(testUser)
	assert.NotNil(t, err)

	//test empty Email
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return true
		},
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, nil
		},
	}

	err = useCase.ChangeUser(&models.User{})
	assert.NotNil(t, err)

	//test internal error ReplaceUser
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return true
		},
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, nil
		},
		ReplaceFunc: func(ID uint64, newUser *models.User) error {
			return errors.New("Internal error")
		},
	}

	err = useCase.ChangeUser(testUser)
	assert.NotNil(t, err)

	//test success
	useCase.repository = &repository.UserRepoMock{
		ContainsFunc: func(user models.User) bool {
			return true
		},
		GetUserByIDFunc: func(ID uint64) (user models.User, e error) {
			return models.User{}, nil
		},
		ReplaceFunc: func(ID uint64, newUser *models.User) error {
			return nil
		},
	}

	err = useCase.ChangeUser(testUser)
	assert.Nil(t, err)

}

func TestUsersUseCase_FindUsers(t *testing.T) {
	testUsername := "mem"

	//test internal error
	useCase.repository = &repository.UserRepoMock{
		GetUsersFunc: func() (users models.Users, e error) {
			return models.Users{}, errors.New("Internal error")
		},
	}

	_, err := useCase.FindUsers(testUsername)
	assert.NotNil(t, err)

	//test ok
	useCase.repository = &repository.UserRepoMock{
		GetUsersFunc: func() (users models.Users, e error) {
			return models.Users{}, nil
		},
	}

	_, err = useCase.FindUsers(testUsername)
	assert.Nil(t, err)
}
