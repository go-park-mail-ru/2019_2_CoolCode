package useCase

import "github.com/AntonPriyma/2019_2_CoolCode/models"

type UsersUseCase interface{
	GetUserByID(id int64) (models.User,error)
	GetUserByEmail(email string) (models.User,error)
	AddUser(user *models.User) error
	ChangeUser(user models.User)
	FindUsers(name string) ([]models.User,error)
}