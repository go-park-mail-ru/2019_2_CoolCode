package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type User struct {
	ID       uint
	Username string
	Email    string
	Name     string
	Password string
	Status   string
	Photo    []byte
}

type UserStore struct {
	users map[uint]*User
	mutex sync.Mutex
	nextID uint
}

type Users struct {
	Users []User `json:"users"`
}

func (userStore *UserStore) readUsers(users Users) {
	for _, user := range users.Users {
		userStore.AddUser(&user)
	}
}

func (userStore *UserStore) saveUsers() {
	usersSlice:=userStore.GetUsers()
	os.Remove("users.txt")
	file, _ := os.Create("users.txt")
	encoder := json.NewEncoder(file)
	encoder.Encode(usersSlice)
}



func NewUserStore() UserStore {
	return UserStore{
		mutex: sync.Mutex{},
		users: make(map[uint]*User, 0),
	}
}

func (userStore UserStore)Contains(user User)  bool{

	for _,v:=range userStore.users{
		if user.Email==v.Email{
			return true
		}
	}

	return false
}

func (userStore UserStore) GetUserByEmail(email string)  (User,error){
	var resultUser User
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()
	for _,v:=range userStore.users{
		if email==v.Email{
			return *v,nil
		}
	}

	return resultUser,errors.New("user not contains")
}

func (userStore UserStore) GetUserByID(ID uint)  (User,error){
	var resultUser User
	userStore.mutex.Lock()
	if user,ok:=userStore.users[ID];ok{
		return *user,nil
	}
	userStore.mutex.Unlock()
	return resultUser,errors.New("user not contains")
}



func (userStore UserStore)ChangeUser(user *User){
	defer userStore.saveUsers()
	userStore.mutex.Lock()
	password:=user.Password
	userStore.users[user.ID]=user
	userStore.users[user.ID].Password=password
	userStore.mutex.Unlock()

}






func (userStore *UserStore)AddUser(newUser *User) error{
	defer userStore.saveUsers()
	userStore.mutex.Lock()
	defer userStore.mutex.Unlock()

	if userStore.Contains(*newUser) {
		log.Println("User contains", newUser)
		return NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
	}
	userStore.nextID++;
	newUser.ID = userStore.nextID
	userStore.users[newUser.ID] = newUser


	return nil
}


func (userStore UserStore)GetUsers() Users {
	var usersSlice Users
	for _, user := range userStore.users {
		usersSlice.Users = append(usersSlice.Users, *user)
	}
	return usersSlice
}

func (userStore UserStore)SavePhoto(file multipart.File,id string) error {
	defer file.Close()



	tempFile, err := ioutil.TempFile("photos", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tempFile.Close()


	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.Rename(tempFile.Name(),"photos/"+id+".png")

	if err != nil {
		fmt.Println(err)
		return err
	}
	tempFile.Write(fileBytes)
	return nil
}
func (userStore *UserStore) GetPhoto(id int) (os.File,error) {
	fileName:=strconv.Itoa(id)
	file,err:=os.Open("photos/"+fileName+".png")
	if err!=nil{
		fmt.Println(err)
		return *file,err
	}
	return *file,nil
}



