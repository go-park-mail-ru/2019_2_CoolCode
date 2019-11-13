package users

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
)

type UsersService interface {
	GetUserByID(id uint64) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	SignUp(user *models.User) error
	Login(user models.User) (models.User, error)
	ChangeUser(user *models.User) error
	FindUsers(name string) (models.Users, error)
}

type UsersServiceImpl struct {
	UseCase useCase.UsersUseCase
}

func (u *UsersServiceImpl) GetUserByID(context.Context, *UserID) (*User, error) {
	panic("implement me")
}

func (u *UsersServiceImpl) GetUserByEmail(context.Context, *UserID) (*User, error) {
	panic("implement me")
}

func (u *UsersServiceImpl) SignUp(context.Context, *User) (*Empty, error) {
	panic("implement me")
}

func (u *UsersServiceImpl) Login(context.Context, *User) (*User, error) {
	panic("implement me")
}

func (u *UsersServiceImpl) ChangeUser(context.Context, *User) (*Empty, error) {
	panic("implement me")
}

func (u *UsersServiceImpl) FindUsers(context.Context, *UserName) (*User, error) {
	panic("implement me")
}

func NewGRPCUsersService(useCase useCase.UsersUseCase) UsersServiceServer {
	return &UsersServiceImpl{UseCase: useCase}
}
