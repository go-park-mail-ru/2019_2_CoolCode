package main

import (
	"github.com/AntonPriyma/2019_2_CoolCode/delivery"
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
	api:=delivery.NewHandlers()

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://boiling-chamber-90136.herokuapp.com"}),
		handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)

	r := mux.NewRouter()

	r.HandleFunc("/users", api.SignUp).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", api.EditProfile).Methods("PUT")
	r.HandleFunc("/logout", api.Logout).Methods("DELETE")
	r.HandleFunc("/photos", api.SavePhoto).Methods("POST")
	r.HandleFunc("/photos/{id:[0-9]+}", api.GetPhoto).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", api.GetUser).Methods("GET")
	r.HandleFunc("/users/{name:[((a-z)|(A-Z))0-9_-]+}", api.FindUsers).Methods("GET")
	r.HandleFunc("/users", api.GetUserBySession).Methods("GET") //TODO:Добавить в API
	log.Println("Server started")

	err := http.ListenAndServe(":8080", corsMiddleware(r))
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return
	}
}
