package dal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	calendar calendarpb.CalendarEventsClient
)

const timeout = 1000

// Init инициализация Data Access Layer
func Init(address string) {
	grpcConnection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	calendar = calendarpb.NewCalendarEventsClient(grpcConnection)
}

// Event событие каледаря
type Event struct {
	UUID     uuid.UUID
	Date     time.Time
	Duration time.Time
	Message  string
}

// GetAllEvents return all user events
func GetAllEvents(user string) ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	events, err := calendar.GetAll(ctx, &calendarpb.GetAllRequest{User: user})
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
func AddEvent(user string, event Event) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendar.Add(ctx, &calendarpb.AddRequest{User: user, Event: &calendarpb.Event{Date: date, Uuid: event.UUID.String(), Message: event.Message}})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}

// EditEvent element to storage
func EditEvent(user string, event Event) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}

	result, err := calendar.Edit(ctx, &calendarpb.EditRequest{User: user, Event: &calendarpb.Event{Date: date, Uuid: event.UUID.String(), Message: event.Message}})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}

// RemoveEvent element to storage
func RemoveEvent(user string, uuid uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	result, err := calendar.Remove(ctx, &calendarpb.RemoveRequst{User: user, Uuid: uuid.String()})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}
