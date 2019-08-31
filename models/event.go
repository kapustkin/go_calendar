package models

import (
	"time"

	"github.com/google/uuid"
)

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}
