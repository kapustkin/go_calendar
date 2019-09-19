package storage

import (
	"time"

	"github.com/google/uuid"
)

// Storage interface
type Storage interface {
	GetAllEvents(UserID int32) ([]Event, error)
	AddEvent(event *Event) (bool, error)
	EditEvent(event *Event) (bool, error)
	RemoveEvent(userID int32, uuid uuid.UUID) (bool, error)

	GetEventsForSend(daysBeforeEvent int32) ([]Event, error)
	SetEventAsSended(userID int32, uuid uuid.UUID) (bool, error)
}

// Event событие каледаря
type Event struct {
	UUID       uuid.UUID
	CreateDate time.Time
	EventDate  time.Time
	Message    string
	UserID     int32
	IsSended   bool
}
