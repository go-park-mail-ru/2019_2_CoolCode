package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"reflect"
	"testing"
)

func TestDBUserStore_GetUserByID(t *testing.T) {
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

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := &DBUserStore{
		DB: db,
	}

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

	// Query Error
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

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

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

func TestDBUserStore_GetUserByEmail(t *testing.T) {
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

	mock.
		ExpectQuery("SELECT id, username, email, name, password, status, phone FROM users WHERE").
		WithArgs(elemEmail).
		WillReturnRows(rows)

	repo := &DBUserStore{
		DB: db,
	}

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

	// Query Error
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

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "user1")

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

func TestDBUserStore_PutUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &DBUserStore{
		DB: db,
	}

	testUser := &models.User{
		Username: "test",
		Email:    "test@mail.ru",
		Name:     "Name Lastname",
		Password: "testpass",
		Status:   "",
		Phone:    "89991234567",
	}

	//OK Query
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

	// Query Error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnError(fmt.Errorf("bad query"))

	id, err = repo.PutUser(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Result Error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = repo.PutUser(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// last id error
	mock.
		ExpectExec(`INSERT INTO users`).
		WithArgs(testUser.Username, testUser.Email, testUser.Name, testUser.Password, testUser.Status, testUser.Phone).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = repo.PutUser(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
