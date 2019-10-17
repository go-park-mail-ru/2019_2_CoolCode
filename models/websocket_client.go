package models

import "github.com/gorilla/websocket"

type WebSocketClient struct {
	ws *websocket.Conn
}
