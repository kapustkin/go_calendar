package calendar

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/models"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

type CalendarServer struct {
}

// Get not implemented
func (c *CalendarServer) Get(ctx context.Context, req *calendarpb.GetRequest) (*calendarpb.EventResponse, error) {
	return nil, fmt.Errorf("this method not implemented. Try later")
}

// GetAll возвращает все записи пользователя
func (c *CalendarServer) GetAll(ctx context.Context, req *calendarpb.GetAllRequest) (*calendarpb.AllEventResponse, error) {
	events := storage.GetAllEvents(req.GetUser())
	var grpcResponse []*calendarpb.Event
	for _, v := range events {
		date, err := ptypes.TimestampProto(v.Date)
		if err == nil {
			grpcResponse = append(grpcResponse, &calendarpb.Event{Uuid: v.UUID.String(), Message: v.Message, Date: date})
		}
	}
	return &calendarpb.AllEventResponse{Events: grpcResponse}, nil
}

// Add добавляет новое событие
func (c *CalendarServer) Add(ctx context.Context, req *calendarpb.AddRequest) (*calendarpb.OperationStatusResponse, error) {
	user := req.GetUser()
	event := req.GetEvent()

	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return &calendarpb.OperationStatusResponse{Sucess: false}, err
	}
	date, err := ptypes.Timestamp(event.GetDate())
	if err != nil {
		return &calendarpb.OperationStatusResponse{Sucess: false}, err
	}
	res := storage.AddEvent(user, models.Event{Date: date, UUID: uuid, Message: event.Message})
	return &calendarpb.OperationStatusResponse{Sucess: res}, nil
}

// Edit редактирует событие
func (c *CalendarServer) Edit(ctx context.Context, req *calendarpb.EditRequest) (*calendarpb.OperationStatusResponse, error) {
	user := req.GetUser()
	event := req.GetEvent()

	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return &calendarpb.OperationStatusResponse{Sucess: false}, err
	}
	date, err := ptypes.Timestamp(event.GetDate())
	if err != nil {
		return &calendarpb.OperationStatusResponse{Sucess: false}, err
	}
	res := storage.EditEvent(user, models.Event{Date: date, UUID: uuid, Message: event.Message})

	return &calendarpb.OperationStatusResponse{Sucess: res}, nil
}

// Remove удаляет событие
func (c *CalendarServer) Remove(ctx context.Context, req *calendarpb.RemoveRequst) (*calendarpb.OperationStatusResponse, error) {
	user := req.GetUser()
	uuidString := req.GetUuid()

	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		return &calendarpb.OperationStatusResponse{Sucess: false}, err
	}
	res := storage.RemoveEvent(user, uuid)
	return &calendarpb.OperationStatusResponse{Sucess: res}, nil
}
