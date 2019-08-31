package calendar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kapustkin/go_calendar/internal/storage"
	"github.com/kapustkin/go_calendar/models"
)

const userFieldName string = "user"
const errReadBody string = "Error read body"
const errParsing string = "Error parsing payload"
const errNotFound string = "Event not found"

// AddEvent for user
func AddEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var data models.Event
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(res, errParsing, http.StatusForbidden)
	}
	uuid, err := uuid.NewUUID()
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	event := models.Event{UUID: uuid, Date: time.Now(), Message: data.Message}

	user := chi.URLParam(req, userFieldName)
	storage.AddEvent(user, event)
}

// EditEvent for user
func EditEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var data models.Event
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(res, errParsing, http.StatusForbidden)
	}

	user := chi.URLParam(req, userFieldName)
	if !storage.EditEvent(user, data) {
		http.Error(res, errNotFound, http.StatusNotFound)
	}
}

// RemoveEvent for user
func RemoveEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var data models.Event
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(res, errParsing, http.StatusNotImplemented)
	}

	user := chi.URLParam(req, userFieldName)
	if !storage.RemoveEvent(user, data) {
		http.Error(res, errNotFound, http.StatusNotFound)
	}
}

// GetEvents all events for user
func GetEvents(res http.ResponseWriter, req *http.Request) {
	user := chi.URLParam(req, userFieldName)
	events := storage.GetAllEvents(user)
	data, err := json.Marshal(events)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	res.Write(data)
}
