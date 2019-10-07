package models

type Channel struct {
	ID string
	Name string
	TotalMSGCount int64
	Members []int
}
