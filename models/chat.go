package models

type Chat struct{
	ID uint64
	Name string
	TotalMSGCount int64
	Members []uint64
}
