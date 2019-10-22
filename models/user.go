package models

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"fullname"`
	Password string `json:"password"`
	Status   string `json:"fstatus"`
	Phone    string `json:"phone"`
}
