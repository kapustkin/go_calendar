package inmemory

import (
	"sync"

	"github.com/google/uuid"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

// DB структура хранилища
type DB struct {
}

type database struct {
	sync.RWMutex
	data map[string]map[uuid.UUID]s.Event
}

var (
	db database
)

// GetAllEvents return all user events
func (d DB) GetAllEvents(user string) ([]s.Event, error) {
	db.RLock()
	defer db.RUnlock()
	data := []s.Event{}
	for _, value := range db.data[user] {
		data = append(data, value)
	}
	return data, nil
}

// AddEvent element to storage
func (d DB) AddEvent(user string, event s.Event) bool {
	db.Lock()
	defer db.Unlock()

	userRec := db.data[user]
	if userRec == nil {
		db.data = make(map[string]map[uuid.UUID]s.Event)
	}

	if _, ok := db.data[user]; !ok {
		db.data[user] = map[uuid.UUID]s.Event{}
	}

	if _, ok := db.data[user][event.UUID]; !ok {
		db.data[user][event.UUID] = event
		return true
	}
	return false
}

// EditEvent edit event
func (d DB) EditEvent(user string, event s.Event) bool {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[user][event.UUID]; ok {
		rec := db.data[user][event.UUID]
		rec.Message = event.Message
		db.data[user][event.UUID] = rec
		return true
	}

	return false
}

// RemoveEvent remove event
func (d DB) RemoveEvent(user string, uuid uuid.UUID) bool {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[user][uuid]; ok {
		delete(db.data[user], uuid)
		return true
	}

	return false
}
