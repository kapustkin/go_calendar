package calendar

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
)

// EventServer grpc interface realization
type EventServer struct {
}

// GetEventServer grpc interface realization
func GetEventServer() *EventServer {
	return &EventServer{}
}

// Get not implemented
func (c *EventServer) Get(ctx context.Context, req *calendarpb.GetRequest) (*calendarpb.GetResponse, error) {
	return nil, fmt.Errorf("this method not implemented. Try later")
}

// GetAll возвращает все записи пользователя
func (c *EventServer) GetAll(ctx context.Context, req *calendarpb.GetAllRequest) (*calendarpb.GetAllResponse, error) {
	events := storage.GetAllEvents(req.GetUser())
	var grpcResponse []*calendarpb.Event
	for _, v := range events {
		date, err := ptypes.TimestampProto(v.Date)
		if err != nil {
			return nil, err
		}
		grpcResponse = append(grpcResponse, &calendarpb.Event{Uuid: v.UUID.String(), Message: v.Message, Date: date})
	}
	return &calendarpb.GetAllResponse{Events: grpcResponse}, nil
}

// Add добавляет новое событие
func (c *EventServer) Add(ctx context.Context, req *calendarpb.AddRequest) (*calendarpb.AddResponse, error) {
	user := req.GetUser()
	event := req.GetEvent()

	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return &calendarpb.AddResponse{Sucess: false}, err
	}
	date, err := ptypes.Timestamp(event.GetDate())
	if err != nil {
		return &calendarpb.AddResponse{Sucess: false}, err
	}
	res := storage.AddEvent(user, storage.Event{Date: date, UUID: uuid, Message: event.Message})
	return &calendarpb.AddResponse{Sucess: res}, nil
}

// Edit редактирует событие
func (c *EventServer) Edit(ctx context.Context, req *calendarpb.EditRequest) (*calendarpb.EditResponse, error) {
	user := req.GetUser()
	event := req.GetEvent()

	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		return &calendarpb.EditResponse{Sucess: false}, err
	}
	date, err := ptypes.Timestamp(event.GetDate())
	if err != nil {
		return &calendarpb.EditResponse{Sucess: false}, err
	}
	res := storage.EditEvent(user, storage.Event{Date: date, UUID: uuid, Message: event.Message})

	return &calendarpb.EditResponse{Sucess: res}, nil
}

// Remove удаляет событие
func (c *EventServer) Remove(ctx context.Context, req *calendarpb.RemoveRequst) (*calendarpb.RemoveResponse, error) {
	user := req.GetUser()
	uuidString := req.GetUuid()

	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		return &calendarpb.RemoveResponse{Sucess: false}, err
	}
	res := storage.RemoveEvent(user, uuid)
	return &calendarpb.RemoveResponse{Sucess: res}, nil
}
