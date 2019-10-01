package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"fullname"`
	Password string `json:"password"`
	Status   string `json:"fstatus"`
	Phone    string `json:"phone"`
}

type Users struct {
	Users []User `json:"users"`
}

type UserStore struct {
	users  map[uint]*User
	mutex  sync.Mutex
	nextID uint
}

func NewUserStore() UserStore {
	return UserStore{
		mutex: sync.Mutex{},
		users: make(map[uint]*User, 0),
	}
}

func (userStore *UserStore) readUsers(users Users) {
	for _, user := range users.Users {
		err := userStore.AddUser(&user)
		if err != nil {
			log.Println("User adding error:", err.Error())
			return
		}
	}
}

func (userStore *UserStore) saveUsers() {
	usersSlice := userStore.GetUsers()
	err := os.Remove("users.txt")
	if err != nil {
		log.Println(`Removing 'users.txt' error:`, err.Error())
		return
	}
	file, err := os.Create("users.txt")
	if err != nil {
		log.Println(`Creating 'users.txt' error:`, err.Error())
		return
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(usersSlice)
	if err != nil {
		log.Println(`JSON encoding error:`, err.Error())
		return
	}
}

func (userStore *UserStore) Contains(user User) bool {
	for _, v := range userStore.users {
		if user.Email == v.Email {
			return true
		}
	}

	return false
}

func (userStore UserStore) GetUserByEmail(email string) (User, error) {
	var resultUser User
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	for _, v := range userStore.users {
		if email == v.Email {
			return *v, nil
		}
	}

	return resultUser, errors.New("user not contains")
}

func (userStore UserStore) GetUserByID(ID uint) (User, error) {
	var resultUser User
	userStore.mutex.Lock()
	if user, ok := userStore.users[ID]; ok {
		return *user, nil
	}
	userStore.mutex.Unlock()
	return resultUser, errors.New("user not contains")
}

func (userStore UserStore) ChangeUser(user *User) {
	defer userStore.saveUsers()

	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	password := user.Password
	oldPassword := userStore.users[user.ID].Password
	userStore.users[user.ID] = user

	if password != "" {
		userStore.users[user.ID].Password = password
	} else {
		userStore.users[user.ID].Password = oldPassword
	}
}

func (userStore *UserStore) AddUser(newUser *User) error {
	defer userStore.saveUsers()
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	if userStore.Contains(*newUser) {
		log.Println("User contains", newUser)
		return NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
	}
	userStore.nextID++
	newUser.ID = userStore.nextID
	userStore.users[newUser.ID] = newUser

	return nil
}

func (userStore UserStore) GetUsers() Users {
	var usersSlice Users
	for _, user := range userStore.users {
		usersSlice.Users = append(usersSlice.Users, *user)
	}
	return usersSlice
}

func (userStore UserStore) SavePhoto(file multipart.File, id string) (returnErr error) {
	defer func() {
		err := file.Close()

		if err != nil && returnErr == nil {
			log.Printf("An error occurred: %v", err)
			returnErr = err
		}
	}()

	tempFile, err := ioutil.TempFile("photos", "upload-*.png")
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	defer func() {
		err := tempFile.Close()

		if err != nil && returnErr == nil {
			log.Printf("An error occurred: %v", err)
			returnErr = err
		}
	}()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}
	err = os.Rename(tempFile.Name(), "photos/"+id+".png")

	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	return nil
}
func (userStore *UserStore) GetPhoto(id int) (os.File, error) {
	fileName := strconv.Itoa(id)
	file, err := os.Open("photos/" + fileName + ".png")
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return *file, err
	}
	return *file, nil
}
