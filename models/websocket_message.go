package models

type WebsocketMessage struct {
	WebsocketEventType int     `json:"event_type"`
	Body               Message `json:"body"`
}
