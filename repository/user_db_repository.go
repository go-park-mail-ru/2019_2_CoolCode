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

func (userStore *DBUserStore) GetUserByID(ID uint64) (models.User, error) {
	user := &models.User{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE id = $1"
	row := userStore.DB.QueryRow(selectStr, ID)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)
	if err != nil {
		return *user, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

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

func (userStore *DBUserStore) GetUserByEmail(email string) (models.User, error) {
	user := &models.User{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	selectStr := "SELECT id, username, email, name, password, status, phone FROM users WHERE email = $1"
	row := userStore.DB.QueryRow(selectStr, email)

	err := row.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)

	if err != nil {
		return *user, models.NewServerError(err, http.StatusInternalServerError, "Can not get user: "+err.Error())
	}

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

func (userStore *DBUserStore) PutUser(newUser *models.User) (uint64, error) {
	var ID uint64

	insertQuery := `INSERT INTO users (username, email, name, password, status, phone) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	rows := userStore.DB.QueryRow(insertQuery,
		newUser.Username, newUser.Email, newUser.Name, newUser.Password, newUser.Status, newUser.Phone)

	err := rows.Scan(&ID)
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not put user: "+err.Error())
	}

	return ID, nil
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
	selectStr := "SELECT id, username, email, password  FROM users WHERE email = $1 or username = $2"
	row := userStore.DB.QueryRow(selectStr, user.Email, user.Username)

	err := row.Scan(&sourceUser.ID, &sourceUser.Username, &sourceUser.Email, &sourceUser.Password)
	if err != nil {
		return false
	}
	return true
}

func (userStore *DBUserStore) GetUsers() (models.Users, error) {
	userSlice := models.Users{}
	var name sql.NullString
	var status sql.NullString
	var phone sql.NullString
	rows, err := userStore.DB.Query("SELECT id, username, email, name, password, status, phone FROM users")
	if err != nil {
		return userSlice, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
	}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &name, &user.Password, &status, &phone)
		if err != nil {
			return userSlice, models.NewServerError(err, http.StatusInternalServerError, "Can not get all users: "+err.Error())
		}
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

	return userSlice, nil
}

func NewUserDBStore(db *sql.DB) UserRepo {
	return &DBUserStore{
		db,
	}
}
