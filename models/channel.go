package models

type Channel struct {
	WorkspaceID   uint64
	ID            uint64
	Name          string
	TotalMSGCount int64
	Members       []uint64
	Admins        []uint64
	CreatorID     uint64
}
