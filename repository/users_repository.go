package repository

import "github.com/go-park-mail-ru/2019_2_CoolCode/models"

type UserRepo interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(ID uint64) (models.User, error)
	PutUser(newUser *models.User) error
	Replace(ID uint64,newUser *models.User) error
	Contains(user models.User) bool
	GetUsers() models.Users
}
