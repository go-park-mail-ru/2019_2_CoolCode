package models

type Chat struct {
	ID            uint64
	Name          string
	TotalMSGCount int64
	Members       []uint64
	LastMessage   string
}

type ResponseChatsArray struct {
	Chats      []Chat
	Workspaces []Workspace
}

func NewChatModel(Name string, ID1 uint64, ID2 uint64) *Chat {
	return &Chat{
		ID:            0,
		Name:          Name,
		TotalMSGCount: 0,
		Members:       []uint64{ID1, ID2},
		LastMessage:   "",
	}
}
