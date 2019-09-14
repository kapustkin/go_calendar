package dal

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const timeout = 1000

type GrpcDal struct {
	connection *grpc.ClientConn
}

// Init инициализация Data Access Layer
func Init(addr string) *GrpcDal {
	conn, _ := grpc.Dial(addr, grpc.WithInsecure())
	return &GrpcDal{conn}
}

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}

// GetAllEvents return all user events
func (g *GrpcDal) GetAllEvents(userID string) ([]Event, error) {
	userid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	events, err := calendarpb.NewCalendarEventsClient(g.connection).GetAll(ctx,
		&calendarpb.GetAllRequest{UserId: int32(userid)})
	if err != nil {
		// gRPC error proc example
		if status.Convert(err).Code() == 666 {
			return nil, fmt.Errorf(status.Convert(err).Message())
		}
		return nil, err
	}

	result := make([]Event, len(events.Events))
	for _, v := range events.Events {
		uuid, err := uuid.Parse(v.Uuid)
		if err != nil {
			return nil, err
		}
		date, err := ptypes.Timestamp(v.GetDate())
		if err != nil {
			return nil, err
		}
		result = append(result, Event{Date: date, UUID: uuid, Message: v.Message})
	}
	return result, nil
}

// AddEvent element to storage
func (g *GrpcDal) AddEvent(userID string, event Event) (bool, error) {
	userid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendarpb.NewCalendarEventsClient(g.connection).Add(ctx,
		&calendarpb.AddRequest{Event: &calendarpb.Event{
			UserId:  int32(userid),
			Date:    date,
			Uuid:    event.UUID.String(),
			Message: event.Message}})
	if err != nil {
		return false, err
	}

	return result.Success, nil
}

// EditEvent element to storage
func (g *GrpcDal) EditEvent(userID string, event Event) (bool, error) {
	userid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendarpb.NewCalendarEventsClient(g.connection).
		Edit(ctx, &calendarpb.EditRequest{
			Event: &calendarpb.Event{
				UserId:  int32(userid),
				Date:    date,
				Uuid:    event.UUID.String(),
				Message: event.Message,
			}})
	if err != nil {
		return false, err
	}

	return result.Success, nil
}

// RemoveEvent element to storage
func (g *GrpcDal) RemoveEvent(userID string, uuid fmt.Stringer) (bool, error) {
	userid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	result, err := calendarpb.NewCalendarEventsClient(g.connection).Remove(ctx,
		&calendarpb.RemoveRequst{
			UserId: int32(userid),
			Uuid:   uuid.String()})
	if err != nil {
		return false, err
	}

	return result.Success, nil
}
