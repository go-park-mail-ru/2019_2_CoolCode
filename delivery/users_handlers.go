package delivery

import (
	"bufio"
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandlers struct {
	Users    useCase.UsersUseCase
	Photos   repository.PhotoRepository
	Sessions repository.SessionRepository
}

func NewUsersHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository) *UserHandlers {
	return &UserHandlers{
		Users:    users,
		Photos:   repository.NewPhotosArrayRepository("photos/"),
		Sessions: sessions,
	}
}

func (handlers *UserHandlers) sendError(err error, w http.ResponseWriter) {
	httpError, ok := err.(models.HTTPError)
	if !ok {
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := httpError.ResponseHeaders() // GetUserByEmail http status code and headers.
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

func (handlers *UserHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	var newUser models.User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&newUser)
	if err != nil {

		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
		return
	}

	err = handlers.Users.SignUp(&newUser)
	if err != nil {
		handlers.sendError(err, w)
		return
	}
}

func (handlers *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	var loginUser models.User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		log.Println("Json decoding error")
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
	}

	user, err := handlers.Users.Login(loginUser)
	if err != nil {
		handlers.sendError(err, w)
	} else {
		token := uuid.New()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: expiration}
		err := handlers.Sessions.Put(cookie.Value, user.ID)
		if err != nil {
			handlers.sendError(err, w)
		}
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
	}

}

func (handlers *UserHandlers) SavePhoto(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := strconv.Itoa(int(user.ID))

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {

		//log.Printf("An error occurred: %v", err)
		w.WriteHeader(500)
		return
	}
	file, _, err := r.FormFile("file")

	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid Photo.")

		handlers.sendError(err, w)
		return
	}

	err = handlers.Photos.SavePhoto(file, id)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Successfully Downloaded File")

}

func (handlers *UserHandlers) GetPhoto(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	file, err := handlers.Photos.GetPhoto(requestedID)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reader := bufio.NewReader(file)
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	bytes := make([]byte, size)
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

func (handlers UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := handlers.parseCookie(sessionID)
	loggedIn := err == nil

	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if requestedID != 0 {
		user, err = handlers.Users.GetUserByID(uint64(requestedID))
	}
	if err != nil {
		log.Println("Get user error")
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid ID.")
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

func (handlers UserHandlers) getUserBySession(w http.ResponseWriter, r *http.Request) {
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

func (handlers *UserHandlers) EditProfile(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

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

	if uint64(requestedID) != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var editUser *models.User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&editUser)
	if editUser.ID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Println("Json decoding error")
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.sendError(err, w)
	}

	handlers.Users.ChangeUser(editUser)
}

func (handlers *UserHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("New request: ", r.Body)

	session, _ := r.Cookie("session_id")
	handlers.Sessions.Remove(session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (handlers UserHandlers) parseCookie(cookie *http.Cookie) (models.User, error) {
	id, err := handlers.Sessions.GetID(cookie.Value)
	user, err := handlers.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, err
	}
}

func (handlers UserHandlers) GetUserBySession(w http.ResponseWriter, r *http.Request) {
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

func (handlers UserHandlers) FindUsers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := handlers.parseCookie(cookie)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if name == "" {
		name = user.Username
	}

	users, err := handlers.Users.FindUsers(name)
	if err != nil {
		handlers.sendError(err, w)
	}
	response, err := json.Marshal(users)
	w.Write(response)

}
