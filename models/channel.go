package models

type Channel struct {
	ID            uint64
	Name          string
	TotalMSGCount int64
	Members       []uint64
	Admins        []uint64
	WorkspaceID   uint64
	CreatorID     uint64
}
