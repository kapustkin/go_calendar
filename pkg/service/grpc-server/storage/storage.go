package storage

import (
	"sync"

	"github.com/google/uuid"
	"github.com/kapustkin/go_calendar/pkg/models"
)

// Storage структура хранилища
type Storage struct {
	sync.RWMutex
	data map[string]map[uuid.UUID]models.Event
}

var (
	storage Storage
)

// GetAllEvents return all user events
func GetAllEvents(user string) []models.Event {
	storage.RLock()
	defer storage.RUnlock()
	data := []models.Event{}
	for _, value := range storage.data[user] {
		data = append(data, value)
	}
	return data
}

// AddEvent element to storage
func AddEvent(user string, event models.Event) bool {
	storage.Lock()
	defer storage.Unlock()

	userRec := storage.data[user]
	if userRec == nil {
		storage.data = make(map[string]map[uuid.UUID]models.Event)
	}

	if _, ok := storage.data[user]; !ok {
		storage.data[user] = map[uuid.UUID]models.Event{}
	}

	if _, ok := storage.data[user][event.UUID]; !ok {
		storage.data[user][event.UUID] = event
		return true
	}
	return false
}

// EditEvent edit event
func EditEvent(user string, event models.Event) bool {
	storage.Lock()
	defer storage.Unlock()

	if _, ok := storage.data[user][event.UUID]; ok {
		rec := storage.data[user][event.UUID]
		rec.Message = event.Message
		storage.data[user][event.UUID] = rec
		return true
	}

	return false
}

// RemoveEvent remove event
func RemoveEvent(user string, uuid uuid.UUID) bool {
	storage.Lock()
	defer storage.Unlock()

	if _, ok := storage.data[user][uuid]; ok {
		delete(storage.data[user], uuid)
		return true
	}

	return false
}
