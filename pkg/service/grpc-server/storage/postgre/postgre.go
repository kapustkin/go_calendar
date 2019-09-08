package postgre

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	// no-lint
	_ "github.com/lib/pq"

	"github.com/google/uuid"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

var connString string

var schema = `
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name text
);

CREATE TABLE events (
	id SERIAL PRIMARY KEY, 
	user_id int,
	uuid text,
	start timestamp,
	finish timestamp,
    comment text NULL
)`

type userTable struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type eventTable struct {
	UUID    string         `db:"uuid"`
	Start   time.Time      `db:"start"`
	Finish  time.Time      `db:"finish"`
	Comment sql.NullString `db:"comment"`
}

// DB структура хранилища
type DB struct {
}

// Init storage
func (d DB) Init(conn string) {
	connString = conn

	_, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatalln(err)
	}

	/*
		db.MustExec(schema)
		tx := db.MustBegin()
		err = tx.Commit()
		if err != nil {
			fmt.Printf("Init DB complete \n")
			return
		}
		fmt.Printf("New DB init complete \n")
	*/
}

// GetAllEvents return all user events
func (d DB) GetAllEvents(user string) ([]s.Event, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	events := []eventTable{}
	err = db.Select(&events, `SELECT uuid,start,finish,comment FROM events WHERE user_id=(SELECT id FROM users WHERE name=$1)`, user)
	res, err := mapEvent(&events)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AddEvent element to storage
func (d DB) AddEvent(user string, event s.Event) (bool, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return false, err
	}
	currentUser := userTable{}
	err = db.Get(&currentUser, `SELECT * FROM users WHERE name=$1 LIMIT 1`, user)
	if err != nil {
		return false, err
	}
	_, err = db.NamedExec(`INSERT INTO events (user_id, uuid,start,finish,comment) VALUES (:user_id,:uuid,:start,:finish,:comment)`,
		map[string]interface{}{
			"user_id": currentUser.ID,
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
func (d DB) EditEvent(user string, event s.Event) (bool, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return false, err
	}
	currentUser := userTable{}
	err = db.Get(&currentUser, `SELECT * FROM users WHERE name=$1 LIMIT 1`, user)
	if err != nil {
		return false, err
	}

	val, err := db.NamedExec(`UPDATE events SET (start,finish,comment) = (:start,:finish,:comment) WHERE user_id = :user_id AND uuid = :uuid`,
		map[string]interface{}{
			"user_id": currentUser.ID,
			"uuid":    event.UUID.String(),
			"start":   event.Date,
			"finish":  event.Date,
			"comment": event.Message,
		})

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
func (d DB) RemoveEvent(user string, uuid uuid.UUID) (bool, error) {
	return false, fmt.Errorf("Not implemented")
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
