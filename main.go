package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//коды ошибок
//200 - успех
//400 - неправильные данные в запросе(логин/пароль и т.д.)
//401 - клиент не авторизован
//405 - неверный метод
//500 - фатальная ошибка на сервере

type ClientError interface {
	Error() string
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
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
		return nil, fmt.Errorf("error while parsing response body: %v", err)
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

type Handlers struct {
	Users    UserStore
	Sessions map[string]uint
}

func NewHandlers() *Handlers {
	return &Handlers{
		Users:    NewUserStore(),
		Sessions: make(map[string]uint, 0),
	}
}

func (handlers *Handlers) sendError(err error, w http.ResponseWriter) {
	clientError, ok := err.(ClientError)
	if !ok {
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := clientError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := clientError.ResponseHeaders() // GetUserByEmail http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)

	_, err = w.Write(body)

	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
}

func (handlers Handlers) savePhoto(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := strconv.Itoa(int(user.ID))

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	file, _, err := r.FormFile("file")

	if err != nil {
		log.Printf("Error Retrieving the File: %v", err)
		err = NewClientError(err, http.StatusBadRequest, "Bad request : invalid Photo.")
		handlers.sendError(err, w)
		return
	}

	err = handlers.Users.SavePhoto(file, id)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Successfully Downloaded File")

}

func (handlers Handlers) getPhoto(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	file, err := handlers.Users.GetPhoto(requestedID)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reader := bufio.NewReader(&file)
	bytes := make([]byte, 10<<20)
	_, err = reader.Read(bytes)

	w.Header().Set("content-type", "multipart/form-data;boundary=1")

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Println("Successfully Uploaded File")

}

func (handlers Handlers) getUser(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := handlers.parseCookie(sessionID)
	loggedIn := err == nil

	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if requestedID != 0 {
		user, err = handlers.Users.GetUserByID(uint(requestedID))
	}
	if err != nil {
		log.Println("Get user error")
		err = NewClientError(err, http.StatusBadRequest, "Bad request : invalid ID.")
		handlers.sendError(err, w)
	}

	user.Password = ""
	body, err := json.Marshal(user)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Println("Successfully Uploaded File")

}

func (handlers Handlers) getUserBySession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Println("Valid user session")
}

func (handlers *Handlers) signUp(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	var newUser User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Println("Json decoding error")
		err = NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
		return
	}
	if newUser.Name == "" {
		newUser.Name = "John Doe"
	}
	if newUser.Username == "" {
		newUser.Username = "Stereo"
	}
	err = handlers.Users.AddUser(&newUser)
	if err != nil {
		handlers.sendError(err, w)
	}
}

func (handlers *Handlers) editProfile(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestedID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("An error occurred: %v", err)
	}

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if requestedID == 0 {
		requestedID = int(user.ID)
	}

	if uint(requestedID) != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var editUser *User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&editUser)
	if editUser.ID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println("Json decoding error")
		err = NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
	}

	handlers.Users.ChangeUser(editUser)
}

func (handlers *Handlers) login(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	var loginUser User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		log.Println("Json decoding error")
		err = NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
	}

	user, err := handlers.Users.GetUserByEmail(loginUser.Email)
	if err == nil {
		if user.Password == loginUser.Password {
			//write cookie
			token := uuid.New()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: expiration}
			handlers.Sessions[cookie.Value] = user.ID
			user.Password = ""
			body, err := json.Marshal(user)
			if err != nil {
				log.Printf("An error occurred: %v", err)
				w.WriteHeader(500)
				return
			}
			http.SetCookie(w, &cookie)
			w.Header().Set("content-type", "application/json")

			_, err = w.Write(body)
			if err != nil {
				log.Printf("An error occurred: %v", err)
				w.WriteHeader(500)
				return
			}
			return

		} else {
			log.Println("Wrong password", user)
			err = NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
			handlers.sendError(err, w)
			return
		}
	}

	log.Println("Unregistered user", loginUser)
	err = NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")
	handlers.sendError(err, w)

}

func (handlers *Handlers) logout(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("Not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	delete(handlers.Sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (handlers Handlers) parseCookie(cookie *http.Cookie) (User, error) {
	id := handlers.Sessions[cookie.Value]
	user, err := handlers.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, NewClientError(nil, http.StatusUnauthorized, "Bad request: no Cookie :(")
	}
}

func main() {
	reader, _ := os.Open("users.txt")

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Printf("An error occurred: %v", err)
		}
	}()

	var users Users
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&users)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return
	}
	api := NewHandlers()
	api.Users.readUsers(users)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://boiling-chamber-90136.herokuapp.com"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()

	r.HandleFunc("/users", api.signUp).Methods("POST")
	r.HandleFunc("/login", api.login).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", api.editProfile).Methods("PUT")
	r.HandleFunc("/logout", api.logout).Methods("POST")
	r.HandleFunc("/photos", api.savePhoto).Methods("POST")
	r.HandleFunc("/photos/{id:[0-9]+}", api.getPhoto).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", api.getUser).Methods("GET")
	r.HandleFunc("/users", api.getUserBySession).Methods("GET") //TODO:Добавить в API
	log.Println("Server started")

	err = http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return
	}
}
