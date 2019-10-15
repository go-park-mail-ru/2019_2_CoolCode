package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
)

type DBUserStore struct {
	DB *sql.DB
}

func (userStore *DBUserStore) GetUserByEmail(email string) (models.User, error) {
	panic("implement me")
}

func (userStore *DBUserStore) GetUserByID(ID uint64) (models.User, error) {
	panic("implement me")
}

func (userStore *DBUserStore) PutUser(newUser *models.User) error {
	panic("implement me")
}

func (userStore *DBUserStore) Replace(ID uint64, newUser *models.User) error {
	panic("implement me")
}

func (userStore *DBUserStore) Contains(user models.User) bool {
	panic("implement me")
}

func (userStore *DBUserStore) GetUsers() models.Users {
	panic("implement me")
}

func NewUserDBStore() UserRepo {
	return &DBUserStore{}
}