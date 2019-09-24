package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//коды ошибок
//200 - успех
//500 - фатальная ошибка на сервере
//401 - клиент не авторизован
//400 - неправильные данные в запросе(логин/пароль и т.д.)
//405 - неверный метод

type ClientError interface {
	Error() string
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

type User struct {
	ID       uint64
	Username string
	Email    string
	Name     string
	Password string
	Status   string
	Photo    []byte
}

type HTTPError struct {
	Cause  error
	Detail string
	Status int
}

func (e HTTPError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

func (e HTTPError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

func (e HTTPError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func NewClientError(err error, status int, detail string) ClientError {
	return HTTPError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

type Users struct {
	Users []User `json:"users"`
}

type Handlers struct {
	Users map[uint64]User
}

func (handlers *Handlers) readUsers(users Users) {
	for _, user := range users.Users {
		handlers.Users[user.ID] = user
	}
}

func (handlers *Handlers) saveUsers() {
	var usersSlice Users
	for _, v := range handlers.Users {
		usersSlice.Users = append(usersSlice.Users, v)
	}
	os.Remove("users.txt")
	file, _ := os.Create("users.txt")
	encoder := json.NewEncoder(file)
	encoder.Encode(usersSlice)
}

func signUp(body io.ReadCloser, users map[uint64]User) error {
	var newUser User
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Println("Json decoding error")
		return NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
	}
	newUser.ID = uint64(len(users))
	if _, contains := users[newUser.ID]; contains {
		log.Println("User contains", newUser)
		return NewClientError(nil, http.StatusBadRequest, "Bad request : user already contains.")
	}
	users[newUser.ID] = newUser
	return nil
}

func main() {
	reader, _ := os.Open("users.txt")
	defer reader.Close()
	var users Users
	decoder := json.NewDecoder(reader)
	_ = decoder.Decode(&users)
	handler := Handlers{
		Users: make(map[uint64]User, 0),
	}
	handler.readUsers(users)
	http.HandleFunc("/sign_up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("content-type") == "application/json" {
			err := signUp(r.Body, handler.Users)
			if err != nil {
				clientError, ok := err.(ClientError)
				if !ok {
					w.WriteHeader(500) // return 500 Internal Server Error.
					return
				}

				body, err := clientError.ResponseBody() // Try to get response body of ClientError.
				if err != nil {
					log.Printf("An error accured: %v", err)
					w.WriteHeader(500)
					return
				}
				status, headers := clientError.ResponseHeaders() // Get http status code and headers.
				for k, v := range headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(status)
				w.Write(body)
				return
			}
			ID, _ := json.Marshal(&map[string]int{
				"studentID": len(handler.Users),
			})
			log.Println("Response body: ",string(ID))
			w.Write(ID)
			handler.saveUsers()

		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		err := login(r.Body, handler.Users)
		if err != nil {
			clientError, ok := err.(ClientError)
			if ok {
				body, err := clientError.ResponseBody()
				if err != nil {
					log.Printf("An error accured: %v", err)
					w.WriteHeader(500)
					return
				}
				status, headers := clientError.ResponseHeaders()
				for k, v := range headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(status)
				w.Write(body)
			}
			return
		}
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: "abcd", Expires: expiration}
		http.SetCookie(w, &cookie)
	})

	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		loggedIn := (err != http.ErrNoCookie)

		if !loggedIn {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			editProfile(r, handler.Users)
			handler.saveUsers()
			w.WriteHeader(200)
		}

	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			log.Println("Not authorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	})

	http.ListenAndServe(":8080", nil)

}

func editProfile(request *http.Request, users map[uint64]User) error {
	var editUser User
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&editUser)
	if err != nil {
		log.Println("Json decoding error")
		return NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
	}
	if _, contains := users[editUser.ID]; !contains {
		log.Println("User not contains", editUser)
		return NewClientError(nil, http.StatusBadRequest, "Bad request : user not contains.")
	}
	users[editUser.ID] = editUser
	return nil

}

func login(body io.ReadCloser, users map[uint64]User) error {
	var loginUser User
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		log.Println("Json decoding error")
		return NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
	}
	if val, ok := users[loginUser.ID]; ok {
		if val.Password == loginUser.Password {
			return nil
		} else {
			log.Println("Wrong password",val)
			return NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
		}
	}

	log.Println("Unregistered user", loginUser)
	return NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")

}
