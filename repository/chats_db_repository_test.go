package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"testing"
)

// RemoveWorkspace

func TestChatsDBRepository_RemoveWorkspace_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := NewChatsDBRepository(db)

	var elemID uint64 = 1

	mock.
		ExpectExec("DELETE FROM workspaces WHERE").
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rowsAffected, err := repo.RemoveWorkspace(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if rowsAffected != 1 {
		t.Errorf("unexpected rowsAffected count: %v", rowsAffected)
		return
	}
}

func TestChatsDBRepository_RemoveWorkspace_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	mock.
		ExpectExec("DELETE FROM workspaces WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.RemoveWorkspace(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

// RemoveChannel

func TestChatsDBRepository_RemoveChannel_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	repo := &ChatsDBRepository{
		db: db,
	}

	mock.
		ExpectExec("DELETE FROM chats WHERE").
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rowsAffected, err := repo.RemoveChannel(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if rowsAffected != 1 {
		t.Errorf("unexpected rowsAffected count: %v", rowsAffected)
		return
	}
}

func TestChatsDBRepository_RemoveChannel_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	repo := &ChatsDBRepository{
		db: db,
	}

	mock.
		ExpectExec("DELETE FROM chats WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.RemoveChannel(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

// RemoveChat

func TestChatsDBRepository_RemoveChat_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	repo := &ChatsDBRepository{
		db: db,
	}

	mock.
		ExpectExec("DELETE FROM chats WHERE").
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rowsAffected, err := repo.RemoveChat(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if rowsAffected != 1 {
		t.Errorf("unexpected rowsAffected count: %v", rowsAffected)
		return
	}
}

func TestChatsDBRepository_RemoveChat_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var elemID uint64 = 1

	repo := &ChatsDBRepository{
		db: db,
	}

	mock.
		ExpectExec("DELETE FROM chats WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.RemoveChat(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

// PutWorkspace

func TestChatsDBRepository_PutWorkspace_Success(t *testing.T) {
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
		Members:   []uint64{1, 2},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false).
		WillReturnResult(sqlmock.NewResult(0, 2))
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
}

func TestChatsDBRepository_PutWorkspace_BeginConnError(t *testing.T) {
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
		Members:   []uint64{1, 2},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

	mock.ExpectBegin().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutWorkspace(testWorkspace)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutWorkspace_FirstQueryConnError(t *testing.T) {
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
		Members:   []uint64{1, 2},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

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
}

func TestChatsDBRepository_PutWorkspace_SecondQueryConnError(t *testing.T) {
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
		Members:   []uint64{1, 2},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false).
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

func TestChatsDBRepository_PutWorkspace_CommitConnError(t *testing.T) {
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
		Members:   []uint64{1, 2},
		Admins:    []uint64{1},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO workspaces`).
		WithArgs(testWorkspace.Name, testWorkspace.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutWorkspace(testWorkspace)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// PutChannel

func TestChatsDBRepository_PutChannel_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChannel := &models.Channel{
		Name:          "TestChannel",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
		Admins:        []uint64{1},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(true, testChannel.TotalMSGCount, testChannel.Name, testChannel.WorkspaceID, testChannel.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	id, err := repo.PutChannel(testChannel)
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

func TestChatsDBRepository_PutChannel_BeginConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChannel := &models.Channel{
		Name:          "TestChannel",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
		Admins:        []uint64{1},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutChannel(testChannel)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChannel_FirstQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChannel := &models.Channel{
		Name:          "TestChannel",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
		Admins:        []uint64{1},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO chats`).
		WithArgs(true, testChannel.TotalMSGCount, testChannel.Name, testChannel.WorkspaceID, testChannel.CreatorID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutChannel(testChannel)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChannel_SecondQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChannel := &models.Channel{
		Name:          "TestChannel",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
		Admins:        []uint64{1},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(true, testChannel.TotalMSGCount, testChannel.Name, testChannel.WorkspaceID, testChannel.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutChannel(testChannel)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChannel_CommitConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChannel := &models.Channel{
		Name:          "TestChannel",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
		Admins:        []uint64{1},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(true, testChannel.TotalMSGCount, testChannel.Name, testChannel.WorkspaceID, testChannel.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutChannel(testChannel)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// PutChat

func TestChatsDBRepository_PutChat_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChat := &models.Chat{
		Name:          "TestChat",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(false, testChat.TotalMSGCount, testChat.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	id, err := repo.PutChat(testChat)
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

func TestChatsDBRepository_PutChat_BeginConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChat := &models.Chat{
		Name:          "TestChat",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
	}

	mock.ExpectBegin().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutChat(testChat)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChat_FirstQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChat := &models.Chat{
		Name:          "TestChat",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO chats`).
		WithArgs(false, testChat.TotalMSGCount, testChat.Name).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutChat(testChat)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChat_SecondQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChat := &models.Chat{
		Name:          "TestChat",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(false, testChat.TotalMSGCount, testChat.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, 1, 2).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	_, err = repo.PutChat(testChat)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_PutChat_CommitConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	testChat := &models.Chat{
		Name:          "TestChat",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2},
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO chats`).
		WithArgs(false, testChat.TotalMSGCount, testChat.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, 1, 2).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit().
		WillReturnError(sql.ErrConnDone)

	_, err = repo.PutChat(testChat)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// UpdateWorkspace

func TestChatsDBRepository_UpdateWorkspace_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE workspaces`).
		WithArgs(testWorkspaceUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM workspaces_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateWorkspace_BeginConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin().
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateWorkspace_FirstQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE workspaces`).
		WithArgs(testWorkspaceUpdated.Name, elemID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateWorkspace_SecondQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE workspaces`).
		WithArgs(testWorkspaceUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM workspaces_users`).
		WithArgs(elemID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateWorkspace_ThirdQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE workspaces`).
		WithArgs(testWorkspaceUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM workspaces_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateWorkspace_CommitConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testWorkspaceUpdated := &models.Workspace{
		ID:        elemID,
		Name:      "TestWorkspaceUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE workspaces`).
		WithArgs(testWorkspaceUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM workspaces_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO workspaces_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit().
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateWorkspace(testWorkspaceUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// UpdateChannels

func TestChatsDBRepository_UpdateChannel_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:            elemID,
		Name:          "testChannelUpdated",
		TotalMSGCount: 5,
		Members:       []uint64{1, 2, 3},
		Admins:        []uint64{1, 3},
		WorkspaceID:   1,
		CreatorID:     1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE chats`).
		WithArgs(testChannelUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM chats_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	err = repo.UpdateChannel(testChannelUpdated)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateChannel_BeginConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:        elemID,
		Name:      "testChannelUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin().
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateChannel(testChannelUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateChannel_FirstQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:        elemID,
		Name:      "testChannelUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE chats`).
		WithArgs(testChannelUpdated.Name, elemID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateChannel(testChannelUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateChannel_SecondQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:        elemID,
		Name:      "testChannelUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE chats`).
		WithArgs(testChannelUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM chats_users`).
		WithArgs(elemID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateChannel(testChannelUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateChannel_ThirdQueryConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:        elemID,
		Name:      "testChannelUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE chats`).
		WithArgs(testChannelUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM chats_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	err = repo.UpdateChannel(testChannelUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestChatsDBRepository_UpdateChannel_CommitConnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &ChatsDBRepository{
		db: db,
	}

	var elemID uint64 = 1

	testChannelUpdated := &models.Channel{
		ID:        elemID,
		Name:      "testChannelUpdated",
		Members:   []uint64{1, 2, 3},
		Admins:    []uint64{1, 3},
		CreatorID: 1,
	}

	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE chats`).
		WithArgs(testChannelUpdated.Name, elemID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM chats_users`).
		WithArgs(elemID).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`INSERT INTO chats_users`).
		WithArgs(1, 1, true, 1, 2, false, 1, 3, true).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit().
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateChannel(testChannelUpdated)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
