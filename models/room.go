package models

type Room struct {
	ID int64
	Channels []*Channel
	Members []int
}
