package models

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WebSocketHub struct {
	Clients          map[string]*WebSocketClient
	AddClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	BroadcastChan    chan []byte
}

func NewHub() *WebSocketHub {
	return &WebSocketHub{
		Clients:          make(map[string]*WebSocketClient),
		AddClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		BroadcastChan:    make(chan []byte),
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case conn := <-h.AddClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.RemoveClient(conn)
		case m := <-h.BroadcastChan:
			h.broadcastMessage(m)
		}
	}
}

func (h *WebSocketHub) RemoveClient(conn *websocket.Conn) {
	delete(h.Clients, conn.RemoteAddr().String())
}
func (h *WebSocketHub) addClient(conn *websocket.Conn) {
	h.Clients[conn.RemoteAddr().String()] = &WebSocketClient{
		ws: conn,
	}
}

func (h *WebSocketHub) broadcastMessage(m []byte) {
	for _, conn := range h.Clients {
		err := conn.ws.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}
}
