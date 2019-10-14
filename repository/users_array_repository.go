package repository

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"log"
	"os"
	"sync"
)

type ArrayUserStore struct {
	users  map[uint64]*models.User
	mutex  sync.Mutex
	nextID uint64
}

func NewArrayUserStore() UserRepo {
	store:=&ArrayUserStore{
		mutex: sync.Mutex{},
		users: make(map[uint64]*models.User, 0),
	}
	reader, _ := os.Open("users.txt")

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Printf("An error occurred: %v", err)
		}
	}()

	var users models.Users
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&users)
	if err != nil {
		log.Printf("An error occurred: %v", err)
	}
	store.readUsers(users)
	return store

}

func (userStore *ArrayUserStore) readUsers(users models.Users) {
	userStore.nextID= uint64(len(users.Users) + 1)
	for _, user := range users.Users {
		userStore.users[user.ID] = user
	}
}

func (userStore *ArrayUserStore) saveUsers() {
	usersSlice := userStore.GetUsers()
	err := os.Remove("users.txt")
	if err != nil {
		log.Println(`Removing 'users.txt' error:`, err.Error())
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

func (userStore *ArrayUserStore) Contains(user models.User) bool {
	for _, v := range userStore.users {
		if user.Email == v.Email {
			return true
		}
	}

	return false
}

func (userStore *ArrayUserStore) GetUsers() models.Users {
	var usersSlice models.Users
	for _, user := range userStore.users {
		usersSlice.Users = append(usersSlice.Users, user)
	}
	return usersSlice
}

func (userStore *ArrayUserStore) GetUserByEmail(email string) (models.User, error) {
	var resultUser models.User
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	for _, v := range userStore.users {
		if email == v.Email {
			return *v, nil
		}
	}

	return resultUser, errors.New("user not contains")
}

func (userStore ArrayUserStore) GetUserByID(ID uint64) (models.User, error) {
	var resultUser models.User
	userStore.mutex.Lock()
	if user, ok := userStore.users[ID]; ok {
		return *user, nil
	}
	userStore.mutex.Unlock()
	return resultUser, errors.New("user not contains")
}

func (userStore *ArrayUserStore) PutUser(newUser *models.User) error {
	defer userStore.saveUsers()
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	if newUser.ID==0 {
		userStore.nextID++
		newUser.ID = userStore.nextID
	}
	userStore.users[newUser.ID] = newUser

	return nil
}

func (userStore *ArrayUserStore) Replace(ID uint64,newUser *models.User) error {
	defer userStore.saveUsers()
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	userStore.users[newUser.ID] = newUser

	return nil
}