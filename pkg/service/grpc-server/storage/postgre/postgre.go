package postgre

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	// no-lint
	_ "github.com/lib/pq"

	"github.com/google/uuid"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

type eventTable struct {
	UUID    string         `db:"uuid"`
	Start   time.Time      `db:"start"`
	Finish  time.Time      `db:"finish"`
	Comment sql.NullString `db:"comment"`
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
func (d *DB) GetAllEvents(userID int32) ([]s.Event, error) {
	events := []eventTable{}
	err := d.db.Select(&events, `SELECT uuid,start,finish,comment FROM events WHERE user_id=$1`, userID)
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
func (d *DB) AddEvent(event *s.Event) (bool, error) {
	_, err := d.db.NamedExec(`
	INSERT INTO events 
	(user_id, uuid,start,finish,comment) VALUES 
	(:user_id,:uuid,:start,:finish,:comment)`,
		map[string]interface{}{
			"user_id": event.UserID,
			"uuid":    event.UUID.String(),
			"start":   event.Date,
			"finish":  event.Date,
			"comment": event.Message,
		})
	if err != nil {
		return false, err
	}
	return true, nil
}

// EditEvent edit event
func (d *DB) EditEvent(event *s.Event) (bool, error) {

	val, err := d.db.NamedExec(`
	UPDATE events SET 
	(start,finish,comment) = 
	(:start,:finish,:comment) 
	WHERE user_id = :user_id AND uuid = :uuid`,
		map[string]interface{}{
			"user_id": event.UserID,
			"uuid":    event.UUID.String(),
			"start":   event.Date,
			"finish":  event.Date,
			"comment": event.Message,
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

func mapEvent(input *[]eventTable) ([]s.Event, error) {
	result := []s.Event{}
	for _, r := range *input {
		uuid, err := uuid.Parse(r.UUID)
		if err != nil {
			return nil, err
		}
		result = append(result, s.Event{
			UUID:    uuid,
			Date:    r.Start,
			Message: r.Comment.String,
		})
	}
	return result, nil
}
