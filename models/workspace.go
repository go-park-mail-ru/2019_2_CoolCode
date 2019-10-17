package models

type Workspace struct {
	ID       uint64
	Channels []*Channel
	Members  []uint64
	Admins   []uint64
	Creator  uint64
}
