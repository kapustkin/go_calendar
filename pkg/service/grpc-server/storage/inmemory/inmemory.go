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
	data map[int32]map[uuid.UUID]s.Event
}

var (
	db database
)

// Init storage
func (d DB) Init(params string) {
	// no need init
}

// GetAllEvents return all user events
func (d DB) GetAllEvents(UserID int32) ([]s.Event, error) {
	db.RLock()
	defer db.RUnlock()
	data := []s.Event{}
	for _, value := range db.data[UserID] {
		data = append(data, value)
	}
	return data, nil
}

// AddEvent element to storage
func (d DB) AddEvent(event *s.Event) (bool, error) {
	db.Lock()
	defer db.Unlock()

	userRec := db.data[event.UserID]
	if userRec == nil {
		db.data = make(map[int32]map[uuid.UUID]s.Event)
	}

	if _, ok := db.data[event.UserID]; !ok {
		db.data[event.UserID] = map[uuid.UUID]s.Event{}
	}

	if _, ok := db.data[event.UserID][event.UUID]; !ok {
		db.data[event.UserID][event.UUID] = *event
		return true, nil
	}
	return false, fmt.Errorf("fail adding record %s", event.UUID)
}

// EditEvent edit event
func (d DB) EditEvent(event *s.Event) (bool, error) {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[event.UserID][event.UUID]; ok {
		rec := db.data[event.UserID][event.UUID]
		rec.Message = event.Message
		db.data[event.UserID][event.UUID] = rec
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", event.UUID)
}

// RemoveEvent remove event
func (d DB) RemoveEvent(userID int32, uuid uuid.UUID) (bool, error) {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.data[userID][uuid]; ok {
		delete(db.data[userID], uuid)
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", uuid)
}
