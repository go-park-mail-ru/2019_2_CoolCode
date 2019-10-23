package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"net/http"
	"strconv"
	"strings"
)

type ChatsDBRepository struct {
	db *sql.DB
}

func (c *ChatsDBRepository) RemoveChat(chatID uint64) error {
	tx, err := c.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not begin transaction in RemoveChat: "+err.Error())
	}
	_, err = tx.Exec("delete from chats_users where chatid=$1", chatID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete users_chats "+
			"in RemoveChat transaction: "+err.Error())
	}

	_, err = tx.Exec("DELETE from chats where id=$1", chatID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete chat in RemoveChat: "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not commit RemoveChat transaction "+err.Error())
	}
	return nil

}

func (c *ChatsDBRepository) GetWorkspaceByID(workspaceID uint64) (models.Workspace, error) {
	var result models.Workspace

	tx, err := c.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not begin transaction for GetWorkspace: "+err.Error())
	}

	row := tx.QueryRow("select id,name,creatorid FROM workspaces where id=$1", workspaceID)

	if err := row.Scan(&result.ID, &result.Name, &result.CreatorID); err != nil {
		return result, models.NewClientError(err, http.StatusBadRequest, "workspace not exists: "+err.Error())
	}

	rows, err := tx.Query("select userid,isadmin from workspaces_users where workspaceid=$1", workspaceID)

	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId for GetWorkspace: "+err.Error())
	}

	for rows.Next() {
		var userID uint64
		var isAdmin bool
		err = rows.Scan(&userID, &isAdmin)
		if err != nil {
			return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId,isAdmin for GetWorkspace: "+err.Error())
		}

		result.Members = append(result.Members, userID)

		if isAdmin {
			result.Admins = append(result.Admins, userID)
		}
	}

	return result, nil
}

func (c *ChatsDBRepository) GetWorkspaces(userID uint64) ([]models.Workspace, error) {
	var result []models.Workspace
	rows, err := c.db.Query("select workspaces_users from workspaces_users where userid=$1", userID)
	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "Can not get workspacesID in GetWorkspaces: "+err.Error())
	}
	if rows == nil {
		return result, nil
	}
	workspacesID := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		if !contains(workspacesID, id) {
			workspacesID = append(workspacesID, id)
		}
	}
	for _, id := range workspacesID {
		workspace, err := c.GetWorkspaceByID(id)
		if err != nil {
			return result, err
		}
		result = append(result, workspace)
	}
	return result, nil
}

func (c *ChatsDBRepository) PutWorkspace(workspace *models.Workspace) (uint64, error) {
	var workspaceID uint64
	tx, err := c.db.Begin()
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not open PutWorkspace transaction "+err.Error())
	}

	defer tx.Rollback()
	_ = c.db.QueryRow("insert into workspaces (name, creatorid) values ($1,$2) returning id", workspace.Name, workspace.CreatorID).Scan(&workspaceID)
	sqlStr := "INSERT into workspaces_users (workspaceid, userid,isAdmin) values "
	var vals []interface{}
	index := 1
	for _, userID := range workspace.Members {
		sqlStr += "($" + strconv.Itoa(index) + "," + "$" + strconv.Itoa(index+1) + "," + "$" + strconv.Itoa(index+2) + "),"
		index += 3
		if contains(workspace.Admins, userID) {
			vals = append(vals, workspaceID, userID, true)
		} else {
			vals = append(vals, workspaceID, userID, false)
		}
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	_, err = c.db.Exec(sqlStr, vals...)
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Put workspace error "+err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not commit PutWorkspace transaction "+err.Error())
	}
	return workspaceID, nil
}

func (c *ChatsDBRepository) PutChannel(channel *models.Channel) (uint64, error) {
	var channelID uint64
	tx, err := c.db.Begin()

	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not open PutChannel transaction "+err.Error())
	}

	defer tx.Rollback()
	_ = c.db.QueryRow("INSERT into chats (ischannel, totalmsgcount, name, workspaceid, creatorid) values ($1,$2,$3,$4,$5) returning id",
		true, channel.TotalMSGCount, channel.Name, channel.WorkspaceID, channel.CreatorID).Scan(&channelID)

	//insert into chats_users
	sqlStr := "INSERT into chats_users (chatid, userid,isAdmin) values "
	var vals []interface{}
	index := 1
	for _, userID := range channel.Members {
		sqlStr += "($" + strconv.Itoa(index) + "," + "$" + strconv.Itoa(index+1) + "," + "$" + strconv.Itoa(index+2) + "),"
		index += 3
		if contains(channel.Admins, userID) {
			vals = append(vals, channelID, userID, true)
		} else {
			vals = append(vals, channelID, userID, false)
		}
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	_, err = c.db.Exec(sqlStr, vals...)

	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Put channel error "+err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not commit PutChannel transaction "+err.Error())
	}
	return channelID, nil
}

func (c *ChatsDBRepository) UpdateWorkspace(workspace *models.Workspace) error {
	tx, err := c.db.Begin()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not begin UpdateWorkspace transaction: "+err.Error())
	}
	defer tx.Rollback()
	_, err = tx.Exec("update workspaces set name = $1 where id=$2", workspace.Name, workspace.ID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not update UpdateWorkspace transaction: "+err.Error())
	}

	_, err = tx.Exec("delete from workspaces_users where workspaceid=$1", workspace.ID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete in UpdateWorkspace transaction: "+err.Error())
	}

	sqlStr := "INSERT into workspaces_users (workspaceid, userid,isAdmin) values "
	var vals []interface{}
	index := 1
	for _, userID := range workspace.Members {
		sqlStr += "($" + strconv.Itoa(index) + "," + "$" + strconv.Itoa(index+1) + "," + "$" + strconv.Itoa(index+2) + "),"
		index += 3
		if contains(workspace.Admins, userID) {
			vals = append(vals, workspace.ID, userID, true)
		} else {
			vals = append(vals, workspace.ID, userID, false)
		}
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	_, err = c.db.Exec(sqlStr, vals...)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not insert workspace_users in "+
			"UpdateWorkspace transaction: "+err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not commit UpdateWorkspace transaction "+err.Error())
	}
	return nil

}

func (c *ChatsDBRepository) GetChannelByID(channelID uint64) (models.Channel, error) {
	var result models.Channel

	tx, err := c.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not begin transaction for GetChannel: "+err.Error())
	}

	row := tx.QueryRow("select id,name,totalmsgcount,creatorid FROM chats where id=$1", channelID)

	if err := row.Scan(&result.ID, &result.Name, &result.TotalMSGCount, &result.CreatorID); err != nil {
		return result, models.NewClientError(err, http.StatusBadRequest, "channel not exists: "+err.Error())
	}

	rows, err := tx.Query("select userid,isadmin from chats_users where chatid=$1", channelID)

	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId for GetChannel: "+err.Error())
	}

	for rows.Next() {
		var userID uint64
		var isAdmin bool
		err = rows.Scan(&userID, &isAdmin)
		if err != nil {
			return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId and isAdmin for GetChannel: "+err.Error())
		}

		result.Members = append(result.Members, userID)

		if isAdmin {
			result.Admins = append(result.Admins, userID)
		}
	}

	return result, nil
}

func (c *ChatsDBRepository) UpdateChannel(channel *models.Channel) error {
	tx, err := c.db.Begin()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not begin UpdateChannel transaction: "+err.Error())
	}
	defer tx.Rollback()
	_, err = tx.Exec("update chats set name = $1 where id=$2", channel.Name, channel.ID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not update UpdateChannel transaction: "+err.Error())
	}

	_, err = tx.Exec("delete from chats_users where chatid=$1", channel.ID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete in UpdateChannel transaction: "+err.Error())
	}

	sqlStr := "INSERT into chats_users (chatid, userid,isAdmin) values "
	var vals []interface{}
	index := 1
	for _, userID := range channel.Members {
		sqlStr += "($" + strconv.Itoa(index) + "," + "$" + strconv.Itoa(index+1) + "," + "$" + strconv.Itoa(index+2) + "),"
		index += 3
		if contains(channel.Admins, userID) {
			vals = append(vals, channel.ID, userID, true)
		} else {
			vals = append(vals, channel.ID, userID, false)
		}
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	_, err = c.db.Exec(sqlStr, vals...)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not insert chats_users in "+
			"UpdateChannel transaction: "+err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not commit UpdateChannel transaction "+err.Error())
	}
	return nil

}

func (c *ChatsDBRepository) RemoveWorkspace(workspaceID uint64) error {
	tx, err := c.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not begin transaction in RemoveWorkspace: "+err.Error())
	}
	_, err = tx.Exec("delete from workspaces_users where workspaceid=$1", workspaceID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete users_workspaces "+
			"in RemoveWorkspace transaction: "+err.Error())
	}
	_, err = tx.Exec("DELETE from workspaces where id=$1", workspaceID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete chat in RemoveWorkspace: "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not commit RemoveWorkspace transaction "+err.Error())
	}
	return nil
}

func (c *ChatsDBRepository) RemoveChannel(channelID uint64) error {
	tx, err := c.db.Begin()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not begin transaction in RemoveChannel: "+err.Error())
	}
	_, err = tx.Exec("DELETE from chats where id=$1", channelID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete chat in RemoveChannel: "+err.Error())
	}
	_, err = tx.Exec("delete from chats_users where chatid=$1", channelID)
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not delete users_chats "+
			"in RemoveChannel transaction: "+err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return models.NewServerError(err, http.StatusInternalServerError, "Can not commit RemoveChannel transaction "+err.Error())
	}
	return nil
}

func (c *ChatsDBRepository) GetChatByID(ID uint64) (models.Chat, error) {
	var result models.Chat

	tx, err := c.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not begin transaction for GetChat: "+err.Error())
	}

	row := tx.QueryRow("select id,name,totalmsgcount FROM chats where id=$1 and ischannel=false", ID)

	if err := row.Scan(&result.ID, &result.Name, &result.TotalMSGCount); err != nil {
		return result, models.NewClientError(err, http.StatusBadRequest, "chat not exists: "+err.Error())
	}

	rows, err := tx.Query("select userid from chats_users where chatid=$1", ID)

	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId for GetChat: "+err.Error())
	}

	for rows.Next() {
		var userID uint64
		err = rows.Scan(&userID)
		if err != nil {
			return result, models.NewServerError(err, http.StatusInternalServerError, "can not get userId for GetChat: "+err.Error())
		}
		result.Members = append(result.Members, userID)
	}

	return result, nil
}

func (c *ChatsDBRepository) PutChat(Chat *models.Chat) (uint64, error) {
	var chatID uint64
	tx, err := c.db.Begin()
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not open PutChat transaction "+err.Error())
	}

	defer tx.Rollback()

	_ = c.db.QueryRow("INSERT into chats (ischannel, totalmsgcount, name) values ($1,$2,$3) returning id",
		false, Chat.TotalMSGCount, Chat.Name).Scan(&chatID)

	//chats_users insert
	sqlStr := "INSERT into chats_users (chatid, userid) values "
	var vals []interface{}
	index := 1
	for _, userID := range Chat.Members {
		sqlStr += "($" + strconv.Itoa(index) + "," + "$" + strconv.Itoa(index+1) + "),"
		index += 2
		vals = append(vals, chatID, userID)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	_, err = c.db.Exec(sqlStr, vals...)
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Put chat error "+err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return 0, models.NewServerError(err, http.StatusInternalServerError, "Can not commit PutChat transaction "+err.Error())
	}
	return chatID, nil

}

func (c *ChatsDBRepository) Contains(Chat models.Chat) error {
	panic("implement me")
}

func (c *ChatsDBRepository) GetChats(userID uint64) ([]models.Chat, error) {
	var result []models.Chat
	rows, err := c.db.Query("select chatid from chats_users where userid=$1", userID)
	if err != nil {
		return result, models.NewServerError(err, http.StatusInternalServerError, "Can not get chatsId in GetChats: "+err.Error())
	}
	if rows == nil {
		return result, nil
	}
	chatsId := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		rows.Scan(&id)
		if !contains(chatsId, id) {
			chatsId = append(chatsId, id)
		}
	}
	for _, id := range chatsId {
		chat, err := c.GetChatByID(id)
		if err == nil {
			result = append(result, chat)
		}
	}
	return result, nil
}

func NewChatsDBRepository(db *sql.DB) ChatsRepository {
	return &ChatsDBRepository{db: db}
}

func contains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
