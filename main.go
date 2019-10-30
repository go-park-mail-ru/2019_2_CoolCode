package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/delivery"
	"github.com/go-park-mail-ru/2019_2_CoolCode/middleware"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
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
	logrus.SetFormatter(&logrus.JSONFormatter{})
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

	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Fatalf("cant connect to redis")
		return
	}

	defer db.Close()
	chatsUseCase := useCase.NewChatsUseCase(repository.NewChatsDBRepository(db))
	messagesUseCase := useCase.NewMessageUseCase(repository.NewMessageDbRepository(db), chatsUseCase)
	usersUseCase := useCase.NewUserUseCase(repository.NewUserDBStore(db))
	notificationsUseCase := useCase.NewNotificationUseCase()
	sessionRepository := repository.NewSessionRedisStore(redisConn)
	usersApi := delivery.NewUsersHandlers(usersUseCase, sessionRepository)
	chatsApi := delivery.NewChatHandlers(usersUseCase, sessionRepository, chatsUseCase)
	notificationApi := delivery.NewNotificationHandlers(usersUseCase, sessionRepository, chatsApi.Chats, notificationsUseCase)
	messagesApi := delivery.NewMessageHandlers(messagesUseCase, usersUseCase, sessionRepository, notificationsUseCase)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()
	handler := middleware.PanicMiddleware(middleware.LogMiddleware(r, logrusLogger))
	r.HandleFunc("/users", usersApi.SignUp).Methods("POST")
	r.HandleFunc("/login", usersApi.Login).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", middleware.AuthMiddleware(usersApi.EditProfile)).Methods("PUT")
	r.Handle("/logout", middleware.AuthMiddleware(usersApi.Logout)).Methods("DELETE")
	r.Handle("/photos", middleware.AuthMiddleware(usersApi.SavePhoto)).Methods("POST")
	r.Handle("/photos/{id:[0-9]+}", middleware.AuthMiddleware(usersApi.GetPhoto)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", middleware.AuthMiddleware(usersApi.GetUser)).Methods("GET")
	r.Handle("/users/{name:[((a-z)|(A-Z))0-9_-]+}", middleware.AuthMiddleware(usersApi.FindUsers)).Methods("GET")
	r.HandleFunc("/users", usersApi.GetUserBySession).Methods("GET") //TODO:Добавить в API

	r.HandleFunc("/chats", chatsApi.PostChat).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/chats", chatsApi.GetChatsByUser).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.GetChatById)).Methods("GET")
	r.Handle("/chats/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.RemoveChat)).Methods("DELETE")

	r.Handle("/channels/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.GetChannelById)).Methods("GET")
	r.Handle("/channels/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.EditChannel)).Methods("PUT")
	r.Handle("/channels/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.RemoveChannel)).Methods("DELETE")
	//TODO: r.Handle("/channels/{id:[0-9]+}/members", middleware.AuthMiddleware(chatsApi.LogoutFromChannel)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}/channels", middleware.AuthMiddleware(chatsApi.PostChannel)).Methods("POST")

	r.Handle("/workspaces/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.GetWorkspaceById)).Methods("GET")
	r.Handle("/workspaces/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.EditWorkspace)).Methods("PUT")
	//TODO: r.Handle("/workspaces/{id:[0-9]+}/members", middleware.AuthMiddleware(chatsApi.LogoutFromWorkspace)).Methods("DELETE")
	r.Handle("/workspaces/{id:[0-9]+}", middleware.AuthMiddleware(chatsApi.RemoveWorkspace)).Methods("DELETE")
	r.Handle("/workspaces", middleware.AuthMiddleware(chatsApi.PostWorkspace)).Methods("POST")
	r.Handle("/chats/{id:[0-9]+}/notifications", middleware.AuthMiddleware(notificationApi.HandleNewWSConnection))

	r.Handle("/chats/{id:[0-9]+}/messages", middleware.AuthMiddleware(messagesApi.SendMessage)).Methods("POST").
		HeadersRegexp("Content-Type", "application/(text|json)")
	r.Handle("/chats/{id:[0-9]+}/messages", middleware.AuthMiddleware(messagesApi.GetMessagesByChatID)).Methods("GET")
	r.Handle("/messages/{id:[0-9]+}", middleware.AuthMiddleware(messagesApi.DeleteMessage)).Methods("DELETE")
	r.Handle("/messages/{id:[0-9]+}", middleware.AuthMiddleware(messagesApi.EditMessage)).Methods("PUT")
	log.Println("Server started")

	err = http.ListenAndServe(":8080", corsMiddleware(handler))
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return
	}
}
