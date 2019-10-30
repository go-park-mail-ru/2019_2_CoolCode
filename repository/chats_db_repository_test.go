package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"testing"
)

func TestChatsDBRepository_PutWorkspace(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testWorkspace := &models.Workspace{
		Name:      "TestWorkspace",
		Members:   []uint64{1},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

	//OK Query
	mock.ExpectBegin()

	mock.
		ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := repo.PutWorkspace(testWorkspace)
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

	// Conn Error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutWorkspace(testWorkspace)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Conn Error
	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutWorkspace(testWorkspace)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
