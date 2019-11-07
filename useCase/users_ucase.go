package useCase

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

//go:generate moq -out users_ucase_mock.go . UsersUseCase
type UsersUseCase interface {
	GetUserByID(id uint64) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	SignUp(user *models.User) error
	Login(user models.User) (models.User, error)
	ChangeUser(user *models.User) error
	FindUsers(name string) (models.Users, error)
}

type usersUseCase struct {
	repository repository.UserRepo
}

func (u *usersUseCase) Login(loginUser models.User) (models.User, error) {
	user, err := u.repository.GetUserByEmail(loginUser.Email)
	if err != nil {
		err = models.NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")
		return user, err
	}

	if comparePasswords(user.Password, loginUser.Password) {
		return user, nil
	} else {
		err = models.NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
		return user, err
	}

}

func NewUserUseCase(repo repository.UserRepo) UsersUseCase {
	return &usersUseCase{
		repository: repo,
	}
}

func (u *usersUseCase) GetUserByID(id uint64) (models.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return user, err
	}
	if !u.Valid(user) {
		return user, models.NewClientError(nil, http.StatusUnauthorized, "Bad request: no such user :(")
	}
	return user, nil
}

func (u *usersUseCase) GetUserByEmail(email string) (models.User, error) {
	return u.repository.GetUserByEmail(email)
}

func (u *usersUseCase) SignUp(newUser *models.User) error {
	if u.repository.Contains(*newUser) {
		return models.NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
	} else {
		if newUser.Name == "" {
			newUser.Name = "John Doe"
		}
		newUser.Password = hashAndSalt(newUser.Password)
		_, err := u.repository.PutUser(newUser)
		if err != nil { // return 500 Internal Server Error.
			return models.NewServerError(err, http.StatusInternalServerError, "")
		}
	}
	return nil
}

func (u *usersUseCase) ChangeUser(user *models.User) error {
	if !u.repository.Contains(*user) {
		return models.NewClientError(nil, http.StatusBadRequest, "Bad request : user not contains.")
	} else {
		oldUser, err := u.repository.GetUserByID(user.ID)
		if err != nil { // return 500 Internal Server Error.
			return models.NewServerError(err, http.StatusInternalServerError, "")
		}
		if user.Email == "" || user.Username == "" {
			return models.NewClientError(nil, 400, "Bad req: empty email or password:(")
		}

		if user.Password == "" {
			user.Password = oldUser.Password
		}
		err = u.repository.Replace(user.ID, user)
		if err != nil { // return 500 Internal Server Error.
			return models.NewServerError(err, http.StatusInternalServerError, "")
		}
	}
	return nil
}

func (u *usersUseCase) FindUsers(username string) (models.Users, error) {
	var result models.Users
	userSlice, err := u.repository.GetUsers()
	if err != nil {
		return result, err
	}
	for _, user := range userSlice.Users {
		if strings.HasPrefix(user.Username, username) {
			user.Password = ""
			result.Users = append(result.Users, user)
		}
	}
	return result, nil
}

func (u *usersUseCase) Valid(user models.User) bool {
	return user.Email != ""
}

func comparePasswords(hashedPassword string, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	if err != nil {
		return false
	}

	return true
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
