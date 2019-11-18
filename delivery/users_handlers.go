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
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	body := r.Body
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&loginUser)
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
		return
	}

	user, err := handlers.Users.Login(loginUser)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
	} else {
		token := uuid.New()
		sessionExpiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_id", Value: token.String(), Expires: sessionExpiration}
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

		//create csrf token
		tokenExpiration := time.Now().Add(24 * time.Hour)
		csrfToken, err := utils.Tokens.Create(user.ID, cookie.Value, tokenExpiration.Unix())
		w.Header().Set("X-CSRF-Token", csrfToken)

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
		handlers.utils.HandleError(err, w, r)
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

func (handlers *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := handlers.parseCookie(sessionID)
	loggedIn := err == nil

	if !loggedIn {
		handlers.utils.HandleError(err, w, r)
		return
	}

	user, err = handlers.Users.GetUserByID(uint64(requestedID))

	if err != nil {
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

func (handlers *UserHandlers) GetUserBySession(w http.ResponseWriter, r *http.Request) {

	sessionID, err := r.Cookie("session_id")
	if err != nil {
		handlers.utils.HandleError(models.NewClientError(err, http.StatusUnauthorized, "Not authorized:("), w, r)
		return
	}

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

func (handlers *UserHandlers) EditProfile(w http.ResponseWriter, r *http.Request) {
	sessionID, _ := r.Cookie("session_id")

	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])

	user, err := handlers.parseCookie(sessionID)
	if err != nil {
		handlers.utils.HandleError(err, w, r)
		return
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
	if err != nil {
		err = models.NewClientError(err, http.StatusBadRequest, "Bad request : invalid JSON.")
		handlers.utils.HandleError(err, w, r)
		return
	}
	if editUser.ID != user.ID {
		err = models.NewClientError(nil, http.StatusUnauthorized,
			fmt.Sprintf("Requested id: %d, user id: %d", editUser.ID, user.ID))
		handlers.utils.HandleError(err, w, r)
		return
	}

	err = handlers.Users.ChangeUser(editUser)

	if err != nil {
		handlers.utils.HandleError(err, w, r)
	}
}

func (handlers *UserHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Cookie("session_id")
	err := handlers.Sessions.Remove(session.Value)
	if err != nil {
		handlers.utils.HandleError(
			models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:("),
			w, r)
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (handlers *UserHandlers) parseCookie(cookie *http.Cookie) (models.User, error) {
	id, err := handlers.Sessions.GetID(cookie.Value)
	if err != nil {
		return models.User{}, models.NewClientError(err, http.StatusUnauthorized, "Bad request : not valid cookie:(")
	}
	user, err := handlers.Users.GetUserByID(id)
	if err == nil {
		return user, nil
	} else {
		return user, err
	}
}

func (handlers *UserHandlers) FindUsers(w http.ResponseWriter, r *http.Request) {
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
