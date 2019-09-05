package calendar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kapustkin/go_calendar/pkg/models"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/dal"
)

const userFieldName string = "user"
const errReadBody string = "Error read body"
const errParsing string = "Error parsing payload"
const errNotFound string = "Event not found"

// GetEvents all events for user
func GetEvents(res http.ResponseWriter, req *http.Request) {
	user := chi.URLParam(req, userFieldName)
	events, err := dal.GetAllEvents(user)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	data, err := json.Marshal(events)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	res.Write(data)
}

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
	result, err := dal.AddEvent(user, event)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if !result {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// EditEvent for user
func EditEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var event models.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		http.Error(res, errParsing, http.StatusForbidden)
	}
	user := chi.URLParam(req, userFieldName)
	result, err := dal.EditEvent(user, event)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if !result {
		http.Error(res, "no record for edit", http.StatusNotFound)
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
	result, err := dal.RemoveEvent(user, data.UUID)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	if !result {
		http.Error(res, "no record for remove", http.StatusNotFound)
	}
}
