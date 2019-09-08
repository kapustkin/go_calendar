package dal

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/dal/pool"
	"google.golang.org/grpc/status"
)

const timeout = 1000

// Init инициализация Data Access Layer
func Init(addr string) {
	pool.Init(addr)
}

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}

// GetAllEvents return all user events
func GetAllEvents(UserID string) ([]Event, error) {
	userid, err := strconv.ParseInt(UserID, 10, 32)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	conn, err := pool.GetConnectionFromPool(&ctx)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	events, err := calendarpb.NewCalendarEventsClient(conn.ClientConn).GetAll(ctx, &calendarpb.GetAllRequest{UserId: int32(userid)})
	if err != nil {
		// gRPC error proc example
		if status.Convert(err).Code() == 666 {
			return nil, fmt.Errorf(status.Convert(err).Message())
		}
		return nil, err
	}

	var result []Event
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
func AddEvent(UserID string, event Event) (bool, error) {
	userid, err := strconv.ParseInt(UserID, 10, 32)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	conn, err := pool.GetConnectionFromPool(&ctx)
	defer conn.Close()
	if err != nil {
		return false, err
	}

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendarpb.NewCalendarEventsClient(conn.ClientConn).Add(ctx, &calendarpb.AddRequest{Event: &calendarpb.Event{UserId: int32(userid), Date: date, Uuid: event.UUID.String(), Message: event.Message}})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}

// EditEvent element to storage
func EditEvent(UserID string, event Event) (bool, error) {
	userid, err := strconv.ParseInt(UserID, 10, 32)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	conn, err := pool.GetConnectionFromPool(&ctx)
	defer conn.Close()
	if err != nil {
		return false, err
	}

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendarpb.NewCalendarEventsClient(conn.ClientConn).
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

	return result.Sucess, nil
}

// RemoveEvent element to storage
func RemoveEvent(UserID string, uuid uuid.UUID) (bool, error) {
	userid, err := strconv.ParseInt(UserID, 10, 32)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	conn, err := pool.GetConnectionFromPool(&ctx)
	defer conn.Close()
	if err != nil {
		return false, err
	}

	result, err := calendarpb.NewCalendarEventsClient(conn.ClientConn).Remove(ctx, &calendarpb.RemoveRequst{UserId: int32(userid), Uuid: uuid.String()})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}
