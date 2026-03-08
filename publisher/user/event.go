package user

import "github.com/google/uuid"

type EventType int

const (
	EventTypeUnknown    EventType = iota
	EventTypeRegistered EventType = iota
)

type Event struct {
	EventID uuid.UUID `json:"event_id"`
	Type    EventType `json:"type"`
	UserID  uuid.UUID `json:"user_id"`
}
