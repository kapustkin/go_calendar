package postgre

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	// no-lint
	_ "github.com/lib/pq"

	"github.com/google/uuid"
	storage "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

type eventTable struct {
	Create    time.Time      `db:"eventcreate"`
	UUID      string         `db:"uuid"`
	Comment   sql.NullString `db:"comment"`
	EventDate time.Time      `db:"eventdate"`
	IsSended  bool           `db:"issended"`
}

// DB структура хранилища
type DB struct {
	db *sqlx.DB
}

// Init storage
func Init(conn string) *DB {
	connection, _ := sqlx.Connect("postgres", conn)
	return &DB{db: connection}
}

// GetAllEvents return all user events
func (d *DB) GetAllEvents(userID int32) ([]storage.Event, error) {
	events := []eventTable{}
	err := d.db.Select(&events, `SELECT uuid,eventcreate,eventdate,comment,issended FROM events WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	res, err := mapEvent(&events)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AddEvent element to storage
func (d *DB) AddEvent(event *storage.Event) (bool, error) {
	_, err := d.db.NamedExec(`
	INSERT INTO events 
	(user_id, uuid,eventcreate,eventdate,comment,issended) VALUES 
	(:user_id,:uuid,:eventcreate,:eventdate,:comment,:issended)`,
		map[string]interface{}{
			"user_id":     event.UserID,
			"uuid":        event.UUID.String(),
			"eventcreate": event.CreateDate,
			"eventdate":   event.EventDate,
			"comment":     event.Message,
			"issended":    event.IsSended,
		})
	if err != nil {
		return false, err
	}
	return true, nil
}

// EditEvent edit event
func (d *DB) EditEvent(event *storage.Event) (bool, error) {

	val, err := d.db.NamedExec(`
	UPDATE events SET 
	(eventdate,comment,issended) = 
	(:eventdate,:comment,:issended) 
	WHERE user_id = :user_id AND uuid = :uuid`,
		map[string]interface{}{
			"user_id":   event.UserID,
			"uuid":      event.UUID.String(),
			"eventdate": event.EventDate,
			"comment":   event.Message,
			"issended":  event.IsSended,
		})
	if err != nil {
		return false, err
	}

	c, err := val.RowsAffected()
	if err != nil {
		return false, err
	}

	if c == 0 {
		return false, fmt.Errorf("record with uuid %s not found", event.UUID)
	}

	return true, nil
}

// RemoveEvent remove event
func (d *DB) RemoveEvent(userID int32, uuid uuid.UUID) (bool, error) {

	val, err := d.db.NamedExec(`DELETE FROM events WHERE user_id = :user_id AND uuid = :uuid`,
		map[string]interface{}{
			"user_id": userID,
			"uuid":    uuid.String(),
		})
	if err != nil {
		return false, err
	}
	c, err := val.RowsAffected()
	if err != nil {
		return false, err
	}

	if c == 0 {
		return false, fmt.Errorf("record with uuid %s not found", uuid)
	}

	return true, nil
}

func (d *DB) GetEventsForSend(daysBeforeEvent int32) ([]storage.Event, error) {
	events := []eventTable{}
	err := d.db.Select(&events,
		`SELECT uuid,eventcreate,eventdate,comment,issended FROM events WHERE eventdate > current_date + interval '%i' day`,
		daysBeforeEvent)
	if err != nil {
		return nil, err
	}
	res, err := mapEvent(&events)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *DB) SetEventAsSended(userID int32, uuid uuid.UUID) (bool, error) {
	val, err := d.db.NamedExec(`
	UPDATE events SET 
	(issended) = (:issended) 
	WHERE user_id = :user_id AND uuid = :uuid`,
		map[string]interface{}{
			"user_id":  userID,
			"uuid":     uuid.String(),
			"issended": true,
		})
	if err != nil {
		return false, err
	}

	c, err := val.RowsAffected()
	if err != nil {
		return false, err
	}

	if c == 0 {
		return false, fmt.Errorf("record with uuid %s not found", uuid)
	}

	return true, nil
}

func mapEvent(input *[]eventTable) ([]storage.Event, error) {
	var result = make([]storage.Event, len(*input))
	for i, r := range *input {
		uuid, err := uuid.Parse(r.UUID)
		if err != nil {
			return nil, err
		}
		result[i] = storage.Event{
			UUID:       uuid,
			CreateDate: r.Create,
			EventDate:  r.EventDate,
			Message:    r.Comment.String,
			IsSended:   r.IsSended,
		}
	}

	//log.Printf("grpc-server.storage.postgre.mapEvent %v", result)
	return result, nil
}
