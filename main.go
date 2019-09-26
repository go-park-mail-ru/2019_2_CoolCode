package _019_2_CoolCode

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
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


type Handlers struct {
	Users UserStore
	Sessions map[string]uint
}

func (handlers *Handlers) sendError(err error,w http.ResponseWriter) {
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
	status, headers := clientError.ResponseHeaders() // GetUserByEmail http status code and headers.
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}



func (handlers *Handlers) signUp(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ",r.Body)
	var newUser User
	body:=r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Println("Json decoding error")
		err =  NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err,w)
	}
	err=handlers.Users.AddUser(&newUser)
	if err!=nil{
		handlers.sendError(err,w)
	}

}

func main() {
	reader, _ := os.Open("users.txt")
	defer reader.Close()
	var users Users
	decoder := json.NewDecoder(reader)
	_ = decoder.Decode(&users)
	handler := Handlers{
		Users:    NewUserStore(),
		Sessions: make(map[string]uint, 0),
	}
	handler.Users.readUsers(users)

	r := mux.NewRouter()
	//r.HandleFunc("/users",addCorsHeader).Methods("OPTIONS")
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Mem"))

	}).Methods("GET")
	r.HandleFunc("/users",handler.signUp).Methods("POST")
	r.HandleFunc("/login",handler.login).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}",handler.editProfile).Methods("PUT")
	r.HandleFunc("/logout",handler.logout).Methods("POST")
	//r.HandleFunc("/users/{id:[0-9]+}").Methods("GET")

	http.ListenAndServe(":8080", r)

}
func (handlers Handlers) parseCookie(cookie *http.Cookie) (User, error) {
	id:=handlers.Sessions[cookie.Value]
	user,err:=handlers.Users.GetUserByID(id)
	if err==nil {
		return user, nil
	} else {
		return user, NewClientError(nil, http.StatusUnauthorized, "Bad request: not Cookie:(")
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
		log.Printf("An error accured: %v", err)
		w.WriteHeader(500)
		return
	}
	user, err := handlers.parseCookie(sessionID)
	loggedIn := err == nil

	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		if uint(requestedID) != user.ID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var editUser *User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&editUser)
		if editUser.ID!=user.ID{
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
}

func (handlers *Handlers) login(w http.ResponseWriter, r *http.Request)  {
	var loginUser User
	body:=r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		log.Println("Json decoding error")
		 err=NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		 handlers.sendError(err,w)
	}

	user,err:=handlers.Users.GetUserByEmail(loginUser.Email)
	if err==nil {
		if user.Password == loginUser.Password {
			//write cookie
			token := uuid.New()
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: expiration}
			handlers.Sessions[cookie.Value] = user.ID
			body, err := json.Marshal(user)
			if err != nil {
				log.Printf("An error accured: %v", err)
				w.WriteHeader(500)
				return
			}
			http.SetCookie(w, &cookie)
			w.Header().Set("content-type","application/json")
			w.Write(body)
			return

		} else {
			log.Println("Wrong password", user)
			err=NewClientError(nil, http.StatusBadRequest, "Bad request: wrong password")
			handlers.sendError(err,w)
		}
	}

	log.Println("Unregistered user", loginUser)
	err = NewClientError(nil, http.StatusBadRequest, "Bad request: malformed data")
	handlers.sendError(err,w)

}

func (handlers *Handlers) logout(w http.ResponseWriter, r *http.Request)  {
		session, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			log.Println("Not authorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
}

func addCorsHeader(w http.ResponseWriter, r *http.Request) {
	log.Println("Handled pre-flight request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}