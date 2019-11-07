package utils

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HandlersUtils struct {
	log *logrus.Logger
}

func NewHandlersUtils(logger *logrus.Logger) HandlersUtils {
	return HandlersUtils{log: logger}
}

func (u *HandlersUtils) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	u.SendError(err, w)
	u.LogError(err, r)
}

func (u *HandlersUtils) SendError(err error, w http.ResponseWriter) {
	httpError, ok := err.(models.HTTPError)
	if !ok {
		w.WriteHeader(500) // return 500 Internal Server Error.
		return
	}

	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		u.log.Error("An error occurred:", err)
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
		u.log.Error("An error occurred:", err)
		w.WriteHeader(500)
		return
	}
}

func (u *HandlersUtils) LogError(err error, r *http.Request) {
	httpError, ok := err.(models.HTTPError)
	if !ok {
		u.log.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"err":         err.Error(),
		}).Error("Internal server error")
		return
	}
	body, err := httpError.ResponseBody() // Try to get response body of ClientError.
	if err != nil {
		u.log.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"err":         err.Error(),
		}).Error("Internal server error")
		return
	}

	u.log.WithFields(logrus.Fields{
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
	}).Error(string(body))

}
