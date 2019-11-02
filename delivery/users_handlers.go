package delivery

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserHandlers struct {
	Users    useCase.UsersUseCase
	Photos   repository.PhotoRepository
	Sessions repository.SessionRepository
	utils    utils.HandlersUtils
}

func NewUsersHandlers(users useCase.UsersUseCase, sessions repository.SessionRepository, utils utils.HandlersUtils) *UserHandlers {
	return &UserHandlers{
		Users:    users,
		Photos:   repository.NewPhotosArrayRepository("photos/"),
		Sessions: sessions,
		utils:    utils,
	}
}

//func (handlers *UserHandlers) sendError(err error, w http.ResponseWriter) {
//	httpError, ok := err.(models.HTTPError)
//	if !ok {
//		w.WriteHeader(500) // return 500 Internal Server Error.
//		return
//	}
//
//	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
//	if err != nil {
//		log.Printf("An error occurred: %v", err)
//		w.WriteHeader(500)
//		return
//	}
//	status, headers := httpError.ResponseHeaders() // GetUserByEmail http status code and headers.
//	for k, v := range headers {
//		w.Header().Set(k, v)
//	}
//	w.WriteHeader(status)
//
//	_, err = w.Write(body)
//
//	if err != nil {
//		log.Printf("An error occurred: %v", err)
//		w.WriteHeader(500)
//		return
//	}
//}
//
//func (hanlers *UserHandlers) logError(err error, r *http.Request) {
//	httpError, ok := err.(models.HTTPError)
//	if !ok {
//		hanlers.log.WithFields(logrus.Fields{
//			"method":      r.Method,
//			"remote_addr": r.RemoteAddr,
//			"err":         err.Error(),
//		}).Error("Internal server error")
//		return
//	}
//	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
//	if err != nil {
//		hanlers.log.WithFields(logrus.Fields{
//			"method":      r.Method,
//			"remote_addr": r.RemoteAddr,
//			"err":         err.Error(),
//		}).Error("Internal server error")
//		return
//	}
//
//	hanlers.log.WithFields(logrus.Fields{
//		"method":      r.Method,
//		"remote_addr": r.RemoteAddr,
//	}).Error(string(body))
//
//}

func (handlers *UserHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&newUser)
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
		return
	}

	err = handlers.Users.SignUp(&newUser)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}
}

func (handlers *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {

	var loginUser models.User
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
	}

	user, err := handlers.Users.Login(loginUser)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
	} else {
		token := uuid.New()
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: expiration}
		err := handlers.Sessions.Put(cookie.Value, user.ID)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
			return
		}
		user.Password = ""
		body, err := json.Marshal(user)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
			return
		}
		http.SetCookie(w, &cookie)
		w.Header().Set("content-type", "application/json")

		_, err = w.Write(body)
		if err != nil {
			handlers.utils.HandleError(err, w, r)
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
		handlers.utils.HandleError(err, w, r)
		return
	}
	file, _, err := r.FormFile("file")

	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid Photo.")
		handlers.utils.HandleError(err, w, r)
		return
	}

	err = handlers.Photos.SavePhoto(file, id)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}
	logrus.WithFields(logrus.Fields{
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
	}).Info("Successfully downloaded file")

}

func (handlers *UserHandlers) GetPhoto(w http.ResponseWriter, r *http.Request) {
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	file, err := handlers.Photos.GetPhoto(requestedID)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
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
		handlers.utils.HandleError(err, w, r)
		return
	}

	logrus.WithFields(logrus.Fields{
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
	}).Info("Successfully uploaded file")

}

func (handlers UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := handlers.parseCookie(sessionID)
	loggedIn := err == nil

	if !loggedIn {
		err = models.NewClientError(nil, http.StatusUnauthorized, "Bad request : unauthorized:(")
		handlers.utils.HandleError(err, w, r)
		return
	}
	if requestedID != 0 {
		user, err = handlers.Users.GetUserByID(uint64(requestedID))
	}
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid ID.")
		handlers.utils.HandleError(err, w, r)
	}

	user.Password = ""
	body, err := json.Marshal(user)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}
}

func (handlers UserHandlers) getUserBySession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		handlers.utils.HandleError(models.NewClientError(nil, http.StatusUnauthorized, "Bad request : unauthorized:("), w, r)
		return
	}
	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		handlers.utils.HandleError(models.NewClientError(nil, http.StatusUnauthorized, "Bad request : unauthorized:("), w, r)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	log.Println("Valid user session")
}

func (handlers *UserHandlers) EditProfile(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	if requestedID == 0 {
		requestedID = int(user.ID)
	}

	if uint64(requestedID) != user.ID {
		err = models.NewClientError(nil, http.StatusUnauthorized,
			fmt.Sprintf("Requested id: %d, user id: %d", requestedID, user.ID))
		handlers.utils.HandleError(err, w, r)
		return
	}

	var editUser *models.User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&editUser)
	if editUser.ID != user.ID {
		err = models.NewClientError(nil, http.StatusUnauthorized,
			fmt.Sprintf("Requested id: %d, user id: %d", editUser.ID, user.ID))
		handlers.utils.HandleError(err, w, r)
		return
	}
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
	}

	err = handlers.Users.ChangeUser(editUser)

	if err != nil {
		handlers.utils.HandleError(err, w, r)
	}
}

func (handlers *UserHandlers) Logout(w http.ResponseWriter, r *http.Request) {

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
		return user, models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
}

func (handlers UserHandlers) GetUserBySession(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}

}

func (handlers UserHandlers) FindUsers(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	cookie, _ := r.Cookie("session_id")

	user, err := handlers.parseCookie(cookie)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	}
	if name == "" {
		name = user.Username
	}

	users, err := handlers.Users.FindUsers(name)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
	}
	response, err := json.Marshal(users)
	w.Write(response)

}
