package storage

import (
	"time"

	"github.com/google/uuid"
)

// Storage interface
type Storage interface {
	GetAllEvents(user string) ([]Event, error)
	AddEvent(user string, event Event) bool
	EditEvent(user string, event Event) bool
	RemoveEvent(user string, uuid uuid.UUID) bool
}

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}
