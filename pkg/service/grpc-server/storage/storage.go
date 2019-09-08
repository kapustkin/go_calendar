package storage

import (
	"time"

	"github.com/google/uuid"
)

// Storage interface
type Storage interface {
	Init(params string)
	GetAllEvents(UserID int32) ([]Event, error)
	AddEvent(event *Event) (bool, error)
	EditEvent(event *Event) (bool, error)
	RemoveEvent(UserID int32, uuid uuid.UUID) (bool, error)
}

// Event событие каледаря
type Event struct {
	UserID   int32
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}
