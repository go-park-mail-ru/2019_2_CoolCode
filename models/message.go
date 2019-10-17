package models

//1 - сообщение
//2 - чувак набирает
type Message struct {
	messageType string
	text        string
	author      User
}
