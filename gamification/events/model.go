package events

import "encoding/json"

type Event struct {
	ID   string `json:"id"  db:"id"`
	Key  string `json:"key" db:"key"`
	Name string `json:"name" db:"name"`
}

type UserEvent struct {
	ID       string          `json:"id" db:"id"`
	EventID  string          `json:"event_id" db:"event_id"`
	UserID   string          `json:"user_id" db:"user_id"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
}
