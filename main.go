package main

import (
	"github.com/go-park-mail-ru/2019_2_CoolCode/delivery"
	"github.com/go-park-mail-ru/2019_2_CoolCode/middleware"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//коды ошибок
//200 - успех
//400 - неправильные данные в запросе(логин/пароль и т.д.)
//401 - клиент не авторизован
//405 - неверный метод
//500 - фатальная ошибка на сервере



func main() {
	userUseCase:=useCase.NewUserUseCase(repository.NewArrayUserStore())
	session:=repository.NewSessionArrayRepository()
	usersApi :=delivery.NewUsersHandlers(userUseCase,session)
	chatsApi :=delivery.NewChatHandlers(userUseCase,session)


	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://boiling-chamber-90136.herokuapp.com"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()
	handler:=middleware.PanicMiddleware(middleware.LogMiddleware(r))
	r.HandleFunc("/users",usersApi.SignUp).Methods("POST")
	r.HandleFunc("/login", usersApi.Login).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", usersApi.EditProfile).Methods("PUT")
	r.HandleFunc("/logout", usersApi.Logout).Methods("DELETE")
	r.HandleFunc("/photos", usersApi.SavePhoto).Methods("POST")
	r.HandleFunc("/photos/{id:[0-9]+}", usersApi.GetPhoto).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", usersApi.GetUser).Methods("GET")
	r.HandleFunc("/users/{name:[((a-z)|(A-Z))0-9_-]+}", usersApi.FindUsers).Methods("GET")
	r.HandleFunc("/users", usersApi.GetUserBySession).Methods("GET") //TODO:Добавить в API

	r.HandleFunc("/chats",chatsApi.PostChat).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/chats",chatsApi.GetChatsByUser).Methods("GET")
	log.Println("Server started")

	err := http.ListenAndServe(":8080", corsMiddleware(handler))
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return
	}
}

//TODO: middleware для ошибок
//TODO: ECHO ???

