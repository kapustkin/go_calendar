package postgre

import (
	"fmt"

	"github.com/google/uuid"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

// DB структура хранилища
type DB struct {
}

// GetAllEvents return all user events
func (d DB) GetAllEvents(user string) ([]s.Event, error) {
	return nil, fmt.Errorf("Not implemented")
}

// AddEvent element to storage
func (d DB) AddEvent(user string, event s.Event) bool {
	return false
}

// EditEvent edit event
func (d DB) EditEvent(user string, event s.Event) bool {
	return false
}

// RemoveEvent remove event
func (d DB) RemoveEvent(user string, uuid uuid.UUID) bool {
	return false
}
