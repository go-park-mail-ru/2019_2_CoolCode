package middleware

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HandlersMiddlwares struct {
	Sessions repository.SessionRepository
	Logger   *logrus.Logger
}

func (m *HandlersMiddlwares) AuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			logrus.SetFormatter(&logrus.TextFormatter{})
			logrus.WithFields(logrus.Fields{
				"method":      r.Method,
				"remote_addr": r.RemoteAddr,
			}).Error("not authorised: request not contains cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//id, err := m.Sessions.GetID(session.Value)
		//if err != nil {
		//	logrus.SetFormatter(&logrus.TextFormatter{})
		//	logrus.WithFields(logrus.Fields{
		//		"method":      r.Method,
		//		"remote_addr": r.RemoteAddr,
		//	}).Error("Cookie not valid" + err.Error())
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}
		//token := r.Header.Get("X-CSRF-Token")
		//_, err = utils.Tokens.Check(id, session.Value, token)
		//if err != nil {
		//	logrus.SetFormatter(&logrus.TextFormatter{})
		//	logrus.WithFields(logrus.Fields{
		//		"method":      r.Method,
		//		"remote_addr": r.RemoteAddr,
		//	}).Error("not authorised: request not contains csrf token")
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}
		////create csrf token
		//tokenExpiration := time.Now().Add(24 * time.Hour)
		//csrfToken, err := utils.Tokens.Create(id, session.Value, tokenExpiration.Unix())
		//w.Header().Set("X-CSRF-Token", csrfToken)
		next(w, r)

	})
}

func (m *HandlersMiddlwares) LogMiddleware(next http.Handler, logrusLogger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.Logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

func (m *HandlersMiddlwares) PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.WithFields(logrus.Fields{
					"method":      r.Method,
					"remote_addr": r.RemoteAddr,
					"error":       err,
				}).Error(r.URL.Path)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
