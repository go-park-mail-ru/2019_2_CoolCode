package models

//1 - сообщение
//2 - чувак набирает

type Message struct {
	ID            uint64 `json:"id"`
	MessageType   int    `json:"type"`
	Text          string `json:"text"`
	AuthorID      uint64 `json:"author_id"`
	ChatID        uint64 `json:"chat_id"`
	FileID        uint64 `json:"file_id"`
	HideForAuthor bool   `json:"-"`
}

type Messages struct {
	Messages []*Message
}
