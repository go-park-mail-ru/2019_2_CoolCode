package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
)

type DBUserStore struct {
	DB *sql.DB
}

func (userStore *DBUserStore) GetUserByEmail(email string) (models.User, error) {
	user := &models.User{}
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE email = $1"
	row := userStore.DB.QueryRow(selectStr, email)

	row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.Status, &user.Phone)
	return *user, nil
}

func (userStore *DBUserStore) GetUserByID(ID uint64) (models.User, error) {
	user := &models.User{}
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE id = $1"
	row := userStore.DB.QueryRow(selectStr, ID)

	row.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.Status, &user.Phone)
	return *user, nil
}

func (userStore *DBUserStore) PutUser(newUser *models.User) error {
	result, _ := userStore.DB.Exec("INSERT INTO users (username, email, name, password, status, phone) VALUES ($1, $2, $3, $4, $5, $6)",
		newUser.Username, newUser.Email, newUser.Name, newUser.Password, newUser.Status, newUser.Phone)
	affected, _ := result.RowsAffected()
	fmt.Printf("Rows affected: %v", affected)
	return nil
}

func (userStore *DBUserStore) Replace(ID uint64, newUser *models.User) error {
	panic("implement me")
}

func (userStore *DBUserStore) Contains(user models.User) bool {
	sourceUser := &models.User{}
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE email = $1"
	row := userStore.DB.QueryRow(selectStr, user.Email)

	err := row.Scan(&sourceUser.ID, &sourceUser.Username, &sourceUser.Email, &sourceUser.Name, &sourceUser.Password, &sourceUser.Status, &sourceUser.Phone)
	if err != nil {
		return false
	}
	return true
}

func (userStore *DBUserStore) GetUsers() models.Users {
	userSlice := models.Users{}

	rows, _ := userStore.DB.Query("SELECT id, username, email, name, password, status, phone FROM users")
	for rows.Next() {
		user := &models.User{}
		rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.Status, &user.Phone)
		userSlice.Users = append(userSlice.Users, user)
	}
	rows.Close()

	return userSlice
}

// TODO: handle errors
func NewUserDBStore() UserRepo {
	connStr := "user=root password=1234 dbname=coolslackdb sslmode=disable"
	db, _ := sql.Open("postgres", connStr)
	db.Ping()

	return &DBUserStore{
		db,
	}
}
