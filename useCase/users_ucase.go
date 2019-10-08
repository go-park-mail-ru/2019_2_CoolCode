package useCase

import (
	"github.com/AntonPriyma/2019_2_CoolCode/models"
	"github.com/AntonPriyma/2019_2_CoolCode/repository"
	"log"
	"net/http"
	"strings"
)

type UsersUseCase interface{
	GetUserByID(id uint64) (models.User,error)
	GetUserByEmail(email string) (models.User,error)
	SignUp(user *models.User) error
	Login(user models.User) (models.User,error)
	ChangeUser(user *models.User) error
	FindUsers(name string) (models.Users,error)
}



type usersUseCase struct {
	repository repository.UserRepo
}

func (u *usersUseCase) Login(loginUser models.User) (models.User, error) {
	user, err := u.repository.GetUserByEmail(loginUser.Email)
	if err!=nil{
		log.Println("Unregistered user", loginUser)
		err = models.NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")
	}

		if user.Password == loginUser.Password {
			return user,nil
		} else {
			log.Println("Wrong password", user)
			err = models.NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
			return user,err
		}


}

func NewUserUseCase(repo repository.UserRepo) UsersUseCase{
	return &usersUseCase{
		repository:repo,
	}
}

func (u *usersUseCase) GetUserByID(id uint64) (models.User, error) {
	return u.repository.GetUserByID(id)
}

func (u *usersUseCase) GetUserByEmail(email string) (models.User, error) {
	return u.repository.GetUserByEmail(email)
}

func(u *usersUseCase) SignUp(newUser *models.User) error {
	if u.repository.Contains(*newUser) {
		log.Println("User contains", newUser)
		return models.NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
	}else {
		if newUser.Name == "" {
			newUser.Name = "John Doe"
		}
		if newUser.Username == "" {
			newUser.Username = "Stereo"
		}
		err := u.repository.PutUser(newUser)
		if err!=nil { // return 500 Internal Server Error.
			log.Printf("An error occurred: %v", err)
			return models.NewServerError(err,http.StatusInternalServerError,"")
		}
	}
	return nil
}

func(u *usersUseCase)ChangeUser(user *models.User) error{
	if !u.repository.Contains(*user){
		return models.NewClientError(nil, http.StatusBadRequest, "Bad request : user not contains.")
	}else {
		err:=u.repository.Replace(user.ID,user)
		if err!=nil { // return 500 Internal Server Error.
			log.Printf("An error occurred: %v", err)
			return models.NewServerError(err,http.StatusInternalServerError,"")
		}
	}
	return nil
}

func(u *usersUseCase) FindUsers(name string) (models.Users,error){
	var result models.Users
	for _,user:=range u.repository.GetUsers().Users{
		if strings.HasPrefix(user.Username,name){
			result.Users = append(result.Users,user)
		}
	}
	return result,nil
}
