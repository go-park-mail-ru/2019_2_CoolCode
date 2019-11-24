package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/delivery"
	"github.com/go-park-mail-ru/2019_2_CoolCode/middleware"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	utils2 "github.com/go-park-mail-ru/2019_2_CoolCode/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
)

//коды ошибок
//200 - успех
//400 - неправильные данные в запросе(логин/пароль и т.д.)
//401 - клиент не авторизован
//405 - неверный метод
//500 - фатальная ошибка на сервере

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1"
	DB_NAME     = "postgres"
)

var (
	redisAddr = flag.String("addr", "redis://localhost:6379", "redis addr")
)

func main() {

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	f, err := os.OpenFile("logs.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrusLogger.Error("Can`t open file:" + err.Error())
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stderr, f)
	logrusLogger.SetOutput(mw)

	utils := utils2.NewHandlersUtils(logrusLogger)
	//init dbConn
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Printf("Error before started: %s", err.Error())
		return
	}
	if db == nil {
		log.Printf("Can not connect to database")
		return
	}

	redisConn := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.DialURL(*redisAddr)
		},
	}
	if err != nil {
		log.Fatalf("cant connect to redis")
		return
	}
	defer redisConn.Close()

	defer db.Close()
	usersRepository := repository.NewUserDBStore(db)
	chatsUseCase := useCase.NewChatsUseCase(repository.NewChatsDBRepository(db), usersRepository)
	messagesUseCase := useCase.NewMessageUseCase(repository.NewMessageDbRepository(db), chatsUseCase)
	usersUseCase := useCase.NewUserUseCase(usersRepository)
	notificationsUseCase := useCase.NewNotificationUseCase()
	sessionRepository := repository.NewSessionRedisStore(redisConn)
	usersApi := delivery.NewUsersHandlers(usersUseCase, sessionRepository, utils)
	chatsApi := delivery.NewChatHandlers(usersUseCase, sessionRepository, chatsUseCase, utils)
	notificationApi := delivery.NewNotificationHandlers(usersUseCase, sessionRepository, chatsApi.Chats, notificationsUseCase, utils)
	messagesApi := delivery.NewMessageHandlers(messagesUseCase, usersUseCase, sessionRepository, notificationsUseCase, utils)
	middlewares := middleware.HandlersMiddlwares{
		Sessions: sessionRepository,
		Logger:   logrusLogger,
	}

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://boiling-chamber-90136.herokuapp.com", "https://boiling-chamber-90136.herokuapp.com", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()
	handler := middlewares.PanicMiddleware(middlewares.LogMiddleware(r, logrusLogger))
	r.HandleFunc("/users", usersApi.SignUp).Methods("POST")
	r.HandleFunc("/login", usersApi.Login).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.EditProfile)).Methods("PUT")
	r.Handle("/logout", middlewares.AuthMiddleware(usersApi.Logout)).Methods("DELETE")
	r.Handle("/photos", middlewares.AuthMiddleware(usersApi.SavePhoto)).Methods("POST")
	r.Handle("/photos/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.GetPhoto)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", middlewares.AuthMiddleware(usersApi.GetUser)).Methods("GET")
	r.Handle("/users/{name:[((a-z)|(A-Z))0-9_-]+}", middlewares.AuthMiddleware(usersApi.FindUsers)).Methods("GET")
	r.HandleFunc("/users", usersApi.GetUserBySession).Methods("GET") //TODO:Добавить в API

	r.HandleFunc("/chats", chatsApi.PostChat).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/chats", chatsApi.GetChatsByUser).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetChatById)).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveChat)).Methods("DELETE")

	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetChannelById)).Methods("GET")
	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.EditChannel)).Methods("PUT")
	r.Handle("/channels/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveChannel)).Methods("DELETE")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.SendMessage)).Methods("POST")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.GetMessagesByChatID)).Methods("GET")
	r.Handle("/channels/{id:[0-9]+}/messages", middlewares.AuthMiddleware(chatsApi.RemoveChannel)).Methods("DELETE")
	//TODO: r.Handle("/channels/{id:[0-9]+}/members", middlewares.AuthMiddleware(chatsApi.LogoutFromChannel)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}/channels", middlewares.AuthMiddleware(chatsApi.PostChannel)).Methods("POST")

	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.GetWorkspaceById)).Methods("GET")
	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.EditWorkspace)).Methods("PUT")
	//TODO: r.Handle("/workspaces/{id:[0-9]+}/members", middlewares.AuthMiddleware(chatsApi.LogoutFromWorkspace)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}", middlewares.AuthMiddleware(chatsApi.RemoveWorkspace)).Methods("DELETE")
	r.Handle("/workspaces", middlewares.AuthMiddleware(chatsApi.PostWorkspace)).Methods("POST")
	r.Handle("/chats/{id:[0-9]+}/notifications", middlewares.AuthMiddleware(notificationApi.HandleNewWSConnection))
	r.Handle("/channels/{id:[0-9]+}/notifications", middlewares.AuthMiddleware(notificationApi.HandleNewWSConnection))

	r.Handle("/chats/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.SendMessage)).Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")
	r.Handle("/chats/{id:[0-9]+}/messages", middlewares.AuthMiddleware(messagesApi.GetMessagesByChatID)).Methods("GET")
	r.Handle("/messages/{text:[((a-z)|(A-Z))0-9_-]+}", middlewares.AuthMiddleware(messagesApi.FindMessages)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", middlewares.AuthMiddleware(messagesApi.DeleteMessage)).Methods("DELETE")
	r.Handle("/messages/{id:[0-9]+}", middlewares.AuthMiddleware(messagesApi.EditMessage)).Methods("PUT")
	r.Handle("/messages/{id:[0-9]+}/likes", middlewares.AuthMiddleware(messagesApi.Like)).Methods("POST")
	log.Println("Server started")
	genetateSSL()

	//err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", corsMiddleware(handler))
	//if err != nil {
	//	logrus.Errorf("Can not listen https, error: %v", err.Error())
	//}

	err = http.ListenAndServe(":8080", corsMiddleware(handler))
	if err != nil {
		logrusLogger.Error(err)
		return
	}
}

func genetateSSL() {
	// Проверяем, доступен ли cert файл.
	err := httpscerts.Check("cert.pem", "key.pem")
	// Если он недоступен, то генерируем новый.
	if err != nil {
		err = httpscerts.Generate("cert.pem", "key.pem", "95.163.209.195:8080")
		if err != nil {
			logrus.Fatal("Ошибка: Не можем сгенерировать https сертификат.")
		}
	}
}
