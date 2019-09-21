package inmemory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	storage "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

// DB структура хранилища
type DB struct {
	db *database
}

type database struct {
	sync.RWMutex
	data map[int32]map[uuid.UUID]storage.Event
}

// Init storage
func Init() *DB {
	storage := make(map[int32]map[uuid.UUID]storage.Event)
	return &DB{db: &database{data: storage}}
}

// GetAllEvents return all user events
func (d *DB) GetAllEvents(userID int32) ([]storage.Event, error) {
	d.db.RLock()
	defer d.db.RUnlock()
	data := []storage.Event{}
	for _, value := range d.db.data[userID] {
		data = append(data, value)
	}
	return data, nil
}

// AddEvent element to storage
func (d *DB) AddEvent(event *storage.Event) (bool, error) {
	d.db.Lock()
	defer d.db.Unlock()

	userRec := d.db.data[event.UserID]
	if userRec == nil {
		d.db.data = make(map[int32]map[uuid.UUID]storage.Event)
	}

	if _, ok := d.db.data[event.UserID]; !ok {
		d.db.data[event.UserID] = map[uuid.UUID]storage.Event{}
	}

	if _, ok := d.db.data[event.UserID][event.UUID]; !ok {
		d.db.data[event.UserID][event.UUID] = *event
		return true, nil
	}
	return false, fmt.Errorf("fail adding record %s", event.UUID)
}

// EditEvent edit event
func (d *DB) EditEvent(event *storage.Event) (bool, error) {
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

func (d *DB) GetEventsForSend() ([]storage.Event, error) {
	//TODO Impl
	return nil, fmt.Errorf("not implemented")
}
func (d *DB) SetEventAsSended(uuid uuid.UUID) (bool, error) {
	//TODO Impl
	return false, fmt.Errorf("not implemented")
}
