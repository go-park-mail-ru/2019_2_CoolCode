package models

type Chat struct{
	ID string
	Name string
	TotalMSGCount int64
	Messages []Message
	Members []int
}
