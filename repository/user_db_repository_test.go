package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"reflect"
	"testing"
)

func TestDBUserStore_GetUserByID_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "name", "password", "status", "phone"})
	expect := []models.User{
		{elemID, "test", "test@mail.ru", "Name Lastname", "testpass", "", "89991234567"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Username, item.Email, item.Name, item.Password, item.Status, item.Phone)
	}

	repo := NewUserDBStore(db)

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	item, err := repo.GetUserByID(elemID)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}
}

func TestDBUserStore_GetUserByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	repo := NewUserDBStore(db)

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetUserByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_GetUserByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

	repo := NewUserDBStore(db)

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.GetUserByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_GetUserByEmail_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemEmail string = "test@mail.ru"

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "name", "password", "status", "phone"})
	expect := []models.User{
		{1, "test", elemEmail, "Name Lastname", "testpass", "", "89991234567"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Username, item.Email, item.Name, item.Password, item.Status, item.Phone)
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemEmail).
		WillReturnRows(rows)

	item, err := repo.GetUserByEmail(elemEmail)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}
}

func TestDBUserStore_GetUserByEmail_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemEmail string = "test@mail.ru"

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemEmail).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetUserByEmail(elemEmail)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_GetUserByEmail_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemEmail string = "test@mail.ru"

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemEmail).
		WillReturnRows(rows)

	_, err = repo.GetUserByEmail(elemEmail)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_PutUser_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	testUser := &models.User{
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name Lastname",
		Password: "testpass",
		Status:   "",
		Phone:    "89991234567",
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.PutUser(testUser)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id != 1 {
		t.Errorf("bad id: want %v, have %v", id, 1)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStore_PutUser_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	testUser := &models.User{
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name Lastname",
		Password: "testpass",
		Status:   "",
		Phone:    "89991234567",
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.PutUser(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStore_Replace_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	testUser := &models.User{
		ID:       elemID,
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name Lastname",
		Password: "testpass",
		Status:   "",
		Phone:    "89991234567",
	}

	testUserChanged := &models.User{
		ID:       elemID,
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name AnotherLastname",
		Password: "AnotherPass",
		Status:   "",
		Phone:    "89991234567",
	}

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "name", "password", "status", "phone"})
	rows = rows.AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone)

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectExec(`UPDATE users`).
		WithArgs(testUserChanged.Username, testUserChanged.Email, testUserChanged.Name, testUserChanged.Password, testUserChanged.Status, testUser.Phone, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Replace(elemID, testUserChanged)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStore_Replace_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	testUser := &models.User{
		ID:       elemID,
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name Lastname",
		Password: "testpass",
		Status:   "",
		Phone:    "89991234567",
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectExec(`UPDATE users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone, elemID).
		WillReturnError(fmt.Errorf("bad query"))

	err = repo.Replace(elemID, testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDBUserStore_Contains_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	testUser := &models.User{
		ID:       1,
		Username: "test",
		Email:    "test@mail.ru",
		Password: "testpass",
	}

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "password"})
	rows = rows.AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Password)

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, password FROM users WHERE").
		WithArgs(testUser.Email, testUser.Username).
		WillReturnRows(rows)

	contains := repo.Contains(*testUser)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !contains {
		t.Errorf("results not match, want %v, have %v", true, contains)
		return
	}
}

func TestDBUserStore_Contains_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	testUser := &models.User{
		ID:       1,
		Username: "test",
		Email:    "test@mail.ru",
		Password: "testpass",
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, password FROM users WHERE").
		WithArgs(testUser.Email, testUser.Username).
		WillReturnError(fmt.Errorf("db_error"))

	contains := repo.Contains(*testUser)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if contains {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_Contains_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	testUser := &models.User{
		ID:       1,
		Username: "test",
		Email:    "test@mail.ru",
		Password: "testpass",
	}

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, password FROM users WHERE").
		WithArgs(testUser.Email, testUser.Username).
		WillReturnRows(rows)

	contains := repo.Contains(*testUser)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if contains {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_GetUsers_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "name", "password", "status", "phone"})
	expect := []models.User{
		{elemID, "test", "test@mail.ru", "Name Lastname", "testpass", "", "89991234567"},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Username, item.Email, item.Name, item.Password, item.Status, item.Phone)
	}

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users").
		WillReturnRows(rows)

	users, err := repo.GetUsers()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(*(users.Users[0]), expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], *(users.Users[0]))
		return
	}
}

func TestDBUserStore_GetUsers_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users").
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetUsers()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDBUserStore_GetUsers_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

	repo := &DBUserStore{
		DB: db,
	}

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users").
		WillReturnRows(rows)

	_, err = repo.GetUsers()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
