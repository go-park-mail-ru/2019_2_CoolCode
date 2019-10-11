package models

import "github.com/gorilla/websocket"

type Client struct {
	id     int
	ws     *websocket.Conn
	ch     chan *Message
	doneCh chan bool
}
