package calendar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/dal"
	logger "github.com/sirupsen/logrus"
)

const userFieldName string = "user"
const errReadBody string = "Error read body"
const errParsing string = "Error parsing payload"

type EventHandler struct {
	dal *dal.GrpcDal
}

// Init calendar event handler
func Init(d *dal.GrpcDal) *EventHandler {
	return &EventHandler{dal: d}
}

// GetEvents all events for user
func (e *EventHandler) GetEvents(res http.ResponseWriter, req *http.Request) {
	user := chi.URLParam(req, userFieldName)
	logger.Infof("userId: %v", user)
	events, err := e.dal.GetAllEvents(user)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Infof("recieved %v events", len(events))
	data, err := json.Marshal(events)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = res.Write(data)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddEvent for user
func (e *EventHandler) AddEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var data dal.Event
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errParsing, http.StatusForbidden)
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	event := dal.Event{UUID: uuid, EventDate: data.EventDate, Message: data.Message}
	user := chi.URLParam(req, userFieldName)
	logger.Infof("event :%v", event)
	result, err := e.dal.AddEvent(user, event)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result {
		message := "insert event failed"
		logger.Warn(message)
		http.Error(res, message, http.StatusInternalServerError)
	}
}

// EditEvent for user
func (e *EventHandler) EditEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var event dal.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errParsing, http.StatusForbidden)
	}
	logger.Infof("event :%v", event)
	user := chi.URLParam(req, userFieldName)
	result, err := e.dal.EditEvent(user, event)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result {
		message := "no record for edit"
		logger.Warn(message)
		http.Error(res, message, http.StatusNotFound)
	}
}

// RemoveEvent for user
func (e *EventHandler) RemoveEvent(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errReadBody, http.StatusForbidden)
	}
	var event dal.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, errParsing, http.StatusNotImplemented)
	}
	logger.Infof("event :%v", event)
	user := chi.URLParam(req, userFieldName)
	result, err := e.dal.RemoveEvent(user, event.UUID)
	if err != nil {
		logger.Errorf(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !result {
		message := "no record for remove"
		logger.Warn(message)
		http.Error(res, message, http.StatusNotFound)
	}
}
