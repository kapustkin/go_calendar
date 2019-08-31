package storage

import (
	"sync"

	"github.com/kapustkin/go_ms_template/models"
)

// Storage структура хранилища
type Storage struct {
	sync.RWMutex
	data map[string][]models.Event
}

var (
	storage Storage
)

// AddEvent element to storage
func AddEvent(user string, element models.Event) {
	storage.Lock()
	defer storage.Unlock()

	rec := storage.data[user]
	if rec == nil {
		storage.data = make(map[string][]models.Event)
	}

	storage.data[user] = append(storage.data[user], element)
}

// EditEvent edit event
func EditEvent(user string, event models.Event) bool {
	storage.Lock()
	defer storage.Unlock()

	for i := 0; i < len(storage.data[user]); i++ {
		if storage.data[user][i].UUID == event.UUID {
			storage.data[user][i].Message = event.Message
			return true
		}
	}
	return false
}

// RemoveEvent remove event
func RemoveEvent(user string, event models.Event) bool {
	storage.Lock()
	defer storage.Unlock()

	for i := 0; i < len(storage.data[user]); i++ {
		if storage.data[user][i].UUID == event.UUID {
			storage.data[user][i], storage.data[user] = storage.data[user][len(storage.data[user])-1], storage.data[user][:len(storage.data[user])-1]
			return true
		}
	}

	return false
}

// GetAllEvents return all user events
func GetAllEvents(user string) []models.Event {
	storage.Lock()
	defer storage.Unlock()
	return storage.data[user]
}
