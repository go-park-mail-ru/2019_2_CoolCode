package models

type Workspace struct {
	ID        uint64
	Name      string
	Channels  []*Channel
	Members   []uint64
	Admins    []uint64
	CreatorID uint64
}
