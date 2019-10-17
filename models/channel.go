package models

type Channel struct {
	ID            uint64
	Name          string
	TotalMSGCount int64
	Members       []uint64
	Admins        []uint64
	Creator       uint64
}
