package dal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/models"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const addr = "localhost:5900"
const timeout = 400

// GetAllEvents return all user events
func GetAllEvents(user string) ([]models.Event, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	calendar := calendarpb.NewCalendarEventsClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	events, err := calendar.GetAll(ctx, &calendarpb.GetAllRequest{User: user})
	if err != nil {
		return nil, err
	}

	var result []models.Event
	fmt.Printf("result recieve: %v", events.Events)
	for _, v := range events.Events {
		uuid, err := uuid.Parse(v.Uuid)
		if err != nil {
			return nil, err
		}
		date, err := ptypes.Timestamp(v.GetDate())
		if err != nil {
			return nil, err
		}
		result = append(result, models.Event{Date: date, UUID: uuid, Message: v.Message})
	}
	return result, nil
}

// AddEvent element to storage
func AddEvent(user string, event models.Event) (bool, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	calendar := calendarpb.NewCalendarEventsClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	date, err := ptypes.TimestampProto(event.Date)
	if err != nil {
		return false, err
	}
	fmt.Printf("1result recieve: %v", event)
	result, err := calendar.Add(ctx, &calendarpb.AddRequest{User: user, Event: &calendarpb.Event{Date: date, Uuid: event.UUID.String(), Message: event.Message}})
	if err != nil {
		return false, err
	}

	return result.Sucess, nil
}
