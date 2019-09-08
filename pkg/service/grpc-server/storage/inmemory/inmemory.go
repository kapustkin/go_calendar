package inmemory

import (
	"fmt"
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

// Init storage
func (d DB) Init(params string) {
	// no need init
}

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
func (d DB) AddEvent(user string, event s.Event) (bool, error) {
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
		return true, nil
	}
	return false, fmt.Errorf("fail adding record %s", event.UUID)
}

// EditEvent edit event
func (d DB) EditEvent(user string, event s.Event) (bool, error) {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[user][event.UUID]; ok {
		rec := db.data[user][event.UUID]
		rec.Message = event.Message
		db.data[user][event.UUID] = rec
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", event.UUID)
}

// RemoveEvent remove event
func (d DB) RemoveEvent(user string, uuid uuid.UUID) (bool, error) {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[user][uuid]; ok {
		delete(db.data[user], uuid)
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", uuid)
}
