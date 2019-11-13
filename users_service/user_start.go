package users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1"
	DB_NAME     = "postgres"
)

func connectDatabase() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return db, err
	}
	if db == nil {
		return db, errors.New("Can not connect to database")
	}
	return db, nil

}

func main() {
	db, err := connectDatabase()
	users := useCase.NewUserUseCase(repository.NewUserDBStore(db))
	service := NewGRPCUsersService(users)

	// Стартуем наш gRPC сервер для прослушивания tcp
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		//
	}
	s := grpc.NewServer()

	// Зарегистрируйте нашу службу на сервере gRPC, это свяжет нашу
	// реализацию с кодом автогенерированного интерфейса для нашего
	// сообщения `Response`, которое мы создали в нашем протобуфе
	RegisterUsersServiceServer(s, service)

	// Регистрация службы ответов на сервере gRPC.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		//
	}
}
