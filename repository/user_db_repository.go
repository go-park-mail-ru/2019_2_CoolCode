package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"net/http"
)

type DBUserStore struct {
	DB *sql.DB
}

func (userStore *DBUserStore) GetUserByEmail(email string) (models.User, error) {
	user := &models.User{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE email = $1"
	row := userStore.DB.QueryRow(selectStr, email)

	row.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)
	if name.Valid {
		user.Name = name.String
	}
	if status.Valid {
		user.Status = status.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
	return *user, nil
}

func (userStore *DBUserStore) GetUserByID(ID uint64) (models.User, error) {
	user := &models.User{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE id = $1"
	row := userStore.DB.QueryRow(selectStr, ID)

	row.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)
	if name.Valid {
		user.Name = name.String
	}
	if status.Valid {
		user.Status = status.String
	}
	if phone.Valid {
		user.Phone = phone.String
	}
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

	_, err := userStore.DB.Exec(
		"UPDATE users SET username = $1, email = $2, name = $3, password = $4, status = $5, phone = $6 WHERE id = $7",
		newUser.Username, newUser.Email, newUser.Name, newUser.Password, newUser.Status, newUser.Phone, ID,
	)

	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not update user: "+err.Error())
	}
	return nil
}

func (userStore *DBUserStore) Contains(user models.User) bool {
	sourceUser := &models.User{}
	selectStr := "SELECT id, username, email, password  FROM users WHERE email = $1"
	row := userStore.DB.QueryRow(selectStr, user.Email)

	err := row.Scan(&sourceUser.ID, &sourceUser.Username, &sourceUser.Email, &sourceUser.Password)
	if err != nil {
		return false
	}
	return true
}

func (userStore *DBUserStore) GetUsers() models.Users {
	userSlice := models.Users{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	rows, _ := userStore.DB.Query("SELECT id, username, email, name, password, status, phone FROM users")
	for rows.Next() {
		user := &models.User{}
		rows.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)
		if name.Valid {
			user.Name = name.String
		}
		if status.Valid {
			user.Status = status.String
		}
		if phone.Valid {
			user.Phone = phone.String
		}
		userSlice.Users = append(userSlice.Users, user)
	}
	rows.Close()

	return userSlice
}

// TODO: handle errors
func NewUserDBStore(db *sql.DB) UserRepo {
	return &DBUserStore{
		db,
	}
}
