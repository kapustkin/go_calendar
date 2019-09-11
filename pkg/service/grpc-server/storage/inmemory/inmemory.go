package inmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

// DB структура хранилища
type DB struct {
	db *database
}

type database struct {
	sync.RWMutex
	data map[int32]map[uuid.UUID]s.Event
}


// Init storage
func Init() *DB {
	storage := make(map[int32]map[uuid.UUID]s.Event)
	return &DB{ db: &database{data: storage}}
}

// GetAllEvents return all user events
func (d *DB) GetAllEvents(UserID int32) ([]s.Event, error) {
	d.db.RLock()
	defer d.db.RUnlock()
	data := []s.Event{}
	for _, value := range d.db.data[UserID] {
		data = append(data, value)
	}
	return data, nil
}

// AddEvent element to storage
func (d *DB) AddEvent(event *s.Event) (bool, error) {
	d.db.Lock()
	defer d.db.Unlock()

	userRec := d.db.data[event.UserID]
	if userRec == nil {
		d.db.data = make(map[int32]map[uuid.UUID]s.Event)
	}

	if _, ok := d.db.data[event.UserID]; !ok {
		d.db.data[event.UserID] = map[uuid.UUID]s.Event{}
	}

	if _, ok := d.db.data[event.UserID][event.UUID]; !ok {
		d.db.data[event.UserID][event.UUID] = *event
		return true, nil
	}
	return false, fmt.Errorf("fail adding record %s", event.UUID)
}

// EditEvent edit event
func (d *DB) EditEvent(event *s.Event) (bool, error) {
	d.db.Lock()
	defer d.db.Unlock()

	if _, ok := d.db.data[event.UserID][event.UUID]; ok {
		rec := d.db.data[event.UserID][event.UUID]
		rec.Message = event.Message
		d.db.data[event.UserID][event.UUID] = rec
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", event.UUID)
}

// RemoveEvent remove event
func (d *DB) RemoveEvent(userID int32, uuid uuid.UUID) (bool, error) {
	d.db.Lock()
	defer d.db.Unlock()

	if _, ok := d.db.data[userID][uuid]; ok {
		delete(d.db.data[userID], uuid)
		return true, nil
	}

	return false, fmt.Errorf("record %s not found", uuid)
}
