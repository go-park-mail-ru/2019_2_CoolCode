package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"reflect"
	"testing"
)

func TestSelectByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "username", "email", "name", "password", "status", "phone"})
	expect := []models.User{
		{elemID, "gdv_fox", "gaaveer@gmail.com", "Daniil Gavrilovsky", "1234", "", "8805553535"},
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


	/*
	// query error
	mock.
		ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.SelectByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "title")

	mock.
		ExpectQuery("SELECT id, title, updated, description FROM items WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	_, err = repo.SelectByID(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}*/

}