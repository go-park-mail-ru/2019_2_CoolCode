package models

import "github.com/gorilla/websocket"

type WebSocketHub struct {
	clients          map[uint64]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

func newHub() *WebSocketHub {
	return &WebSocketHub{
		clients:          make(map[uint64]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}
