package storage

import (
	"sync"
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

// Storage структура хранилища
type Storage struct {
	sync.RWMutex
	data map[string]map[uuid.UUID]Event
}

var (
	storage Storage
)

// GetAllEvents return all user events
func GetAllEvents(user string) []Event {
	storage.RLock()
	defer storage.RUnlock()
	data := []Event{}
	for _, value := range storage.data[user] {
		data = append(data, value)
	}
	return data
}

// AddEvent element to storage
func AddEvent(user string, event Event) bool {
	storage.Lock()
	defer storage.Unlock()

	userRec := storage.data[user]
	if userRec == nil {
		storage.data = make(map[string]map[uuid.UUID]Event)
	}

	if _, ok := storage.data[user]; !ok {
		storage.data[user] = map[uuid.UUID]Event{}
	}

	if _, ok := storage.data[user][event.UUID]; !ok {
		storage.data[user][event.UUID] = event
		return true
	}
	return false
}

// EditEvent edit event
func EditEvent(user string, event Event) bool {
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
