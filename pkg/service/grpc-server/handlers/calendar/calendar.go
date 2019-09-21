package calendar

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/kapustkin/go_calendar/pkg/api/v1"
	storage "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

// EventServer grpc interface realization
type EventServer struct {
	db storage.Storage
}

// GetEventServer grpc interface realization
func GetEventServer(store *storage.Storage) *EventServer {
	return &EventServer{db: *store}
}

// Get not implemented
func (eventServer *EventServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	logger.Errorf("this method not implemented")
	return nil, fmt.Errorf("this method not implemented")
}

// GetAll возвращает все записи пользователя
func (eventServer *EventServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	logger.Infof("get events for userId - %v", req.GetUserId())
	events, err := eventServer.db.GetAllEvents(req.GetUserId())
	if err != nil {
		logger.Errorf(err.Error())
		return nil, status.Error(666, err.Error())
	}
	logger.Infof("%v events received", len(events))
	grpcResponse := make([]*pb.Event, len(events))
	for i, v := range events {
		evDate, err := ptypes.TimestampProto(v.EventDate)
		if err != nil {
			logger.Errorf(err.Error())
			return nil, err
		}
		crDate, err := ptypes.TimestampProto(v.CreateDate)
		if err != nil {
			logger.Errorf(err.Error())
			return nil, err
		}

		grpcResponse[i] = &pb.Event{
			Uuid:       v.UUID.String(),
			Message:    v.Message,
			CreateDate: crDate,
			EventDate:  evDate,
			IsSended:   v.IsSended}
	}

	return &pb.GetAllResponse{Events: grpcResponse}, nil
}

// Add добавляет новое событие
func (eventServer *EventServer) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	event := req.GetEvent()
	logger.Infof("event - %v", event)
	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.AddResponse{Success: false}, err
	}
	createDate, err := ptypes.Timestamp(event.GetCreateDate())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.AddResponse{Success: false}, err
	}
	eventDate, err := ptypes.Timestamp(event.GetEventDate())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.AddResponse{Success: false}, err
	}
	res, err := eventServer.db.AddEvent(&storage.Event{
		UserID:     event.GetUserId(),
		CreateDate: createDate,
		EventDate:  eventDate,
		UUID:       uuid,
		Message:    event.Message,
		IsSended:   false})
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.AddResponse{Success: false}, err
	}
	return &pb.AddResponse{Success: res}, nil
}

// Edit редактирует событие
func (eventServer *EventServer) Edit(ctx context.Context, req *pb.EditRequest) (*pb.EditResponse, error) {
	event := req.GetEvent()
	logger.Infof("event - %v", event)
	uuid, err := uuid.Parse(event.GetUuid())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.EditResponse{Success: false}, err
	}
	evDate, err := ptypes.Timestamp(event.GetEventDate())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.EditResponse{Success: false}, err
	}
	res, err := eventServer.db.EditEvent(
		&storage.Event{
			UUID:      uuid,
			UserID:    event.GetUserId(),
			EventDate: evDate,
			Message:   event.Message,
			IsSended:  event.IsSended})
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.EditResponse{Success: false}, err
	}
	return &pb.EditResponse{Success: res}, nil
}

// Remove удаляет событие
func (eventServer *EventServer) Remove(ctx context.Context, req *pb.RemoveRequst) (*pb.RemoveResponse, error) {
	uuid, err := uuid.Parse(req.GetUuid())
	logger.Infof("event uuid - %v", uuid)
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.RemoveResponse{Success: false}, err
	}

	res, err := eventServer.db.RemoveEvent(req.GetUserId(), uuid)
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.RemoveResponse{Success: false}, err
	}
	return &pb.RemoveResponse{Success: res}, nil
}

// GetEventsForSend получает события для рассылки
func (eventServer *EventServer) GetEventsForSend(ctx context.Context,
	req *pb.GetEventsForSendRequest) (*pb.GetEventsForSendResponse, error) {
	events, err := eventServer.db.GetEventsForSend()
	if err != nil {
		logger.Errorf(err.Error())
		return nil, fmt.Errorf("getEventsForSend failed: %v", err.Error())
	}
	logger.Infof("events count - %v", len(events))
	grpcResponse := make([]*pb.Event, len(events))
	for i, v := range events {
		evDate, err := ptypes.TimestampProto(v.EventDate)
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("event date parse failed: %v", err)
		}
		crDate, err := ptypes.TimestampProto(v.CreateDate)
		if err != nil {
			logger.Errorf(err.Error())
			return nil, fmt.Errorf("create date parse failed: %v", err)
		}

		grpcResponse[i] = &pb.Event{
			Uuid:       v.UUID.String(),
			Message:    v.Message,
			CreateDate: crDate,
			EventDate:  evDate,
			IsSended:   v.IsSended}
	}

	return &pb.GetEventsForSendResponse{Events: grpcResponse}, nil
}

// SetEventAsSended отмечает событие как отправленное
func (eventServer *EventServer) SetEventAsSent(ctx context.Context,
	req *pb.SetEventAsSentRequest) (*pb.SetEventAsSentResponse, error) {
	uuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.SetEventAsSentResponse{Success: false}, fmt.Errorf("parse uuid failed %v", err.Error())
	}
	logger.Infof("event uuid - %v", uuid)
	res, err := eventServer.db.SetEventAsSended(uuid)
	if err != nil {
		logger.Errorf(err.Error())
		return &pb.SetEventAsSentResponse{Success: false}, err
	}
	return &pb.SetEventAsSentResponse{Success: res}, nil
}
