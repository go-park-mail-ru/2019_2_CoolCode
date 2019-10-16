package models

import "database/sql"

type User struct {
	ID       uint64         `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Name     sql.NullString `json:"fullname"`
	Password string         `json:"password"`
	Status   sql.NullString `json:"fstatus"`
	Phone    sql.NullString `json:"phone"`
}
