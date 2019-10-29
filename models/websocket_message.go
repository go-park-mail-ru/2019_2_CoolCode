package models

type WebsocketMessage struct {
	WebsocketEventType int    `json:"event_type"`
	Body               []byte `json:"body"`
}
