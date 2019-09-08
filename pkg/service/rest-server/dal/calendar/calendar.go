package calendar

import (
	"log"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	address        string
	calendar       calendarpb.CalendarEventsClient
	grpcConnection *grpc.ClientConn
)

const timeout = 1000

// Init инициализация Data Access Layer
func Init(addr string) {
	address = addr
}

// GetCalendarService return calendarEventClient
func GetCalendarService() (calendarpb.CalendarEventsClient, error) {
	if grpcConnection == nil || grpcConnection.GetState() != connectivity.Ready {
		grpcConnection, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		calendar = calendarpb.NewCalendarEventsClient(grpcConnection)
	}

	return calendar, nil
}
