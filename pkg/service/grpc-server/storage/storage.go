package storage

import (
	"time"

	"github.com/google/uuid"
)

// Storage interface
type Storage interface {
	Init(params string)
	GetAllEvents(user string) ([]Event, error)
	AddEvent(user string, event Event) (bool, error)
	EditEvent(user string, event Event) (bool, error)
	RemoveEvent(user string, uuid uuid.UUID) (bool, error)
}

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}
