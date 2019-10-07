package repository

import "github.com/AntonPriyma/2019_2_CoolCode/models"

type UserRepo interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(ID uint) (models.User, error)
	PutUser(user models.User) error
	Contains(email string) bool
}
