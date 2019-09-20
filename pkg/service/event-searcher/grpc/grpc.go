package grpc

import (
	"context"
	"time"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/event-searcher/config"
	"google.golang.org/grpc"
)

const timeout = 1000

type Server struct {
	connection *grpc.ClientConn
}

type Event struct {
	UUID    string
	User    string
	Message string
	Date    string
}

func Init(c *config.Config) *Server {
	conn, _ := grpc.Dial(c.GrpcConnection, grpc.WithInsecure())
	return &Server{conn}
}

func (g *Server) GetEventsForSend() ([]Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	events, err := calendarpb.NewCalendarEventsClient(g.connection).GetEventsForSend(ctx,
		&calendarpb.GetEventsForSendRequest{})
	if err != nil {
		return nil, err
	}

	result := make([]Event, len(events.Events))

	for i, v := range events.Events {
		result[i] = Event{
			UUID:    v.GetUuid(),
			Date:    v.GetEventDate().String(),
			User:    v.GetUserName(),
			Message: v.GetMessage()}
	}

	return result, nil
}

func (g *Server) SetEventAsSended(uuid string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
	defer cancel()

	result, err := calendarpb.NewCalendarEventsClient(g.connection).SetEventAsSent(ctx,
		&calendarpb.SetEventAsSentRequest{Uuid: uuid})
	if err != nil {
		return false, err
	}
	return result.Success, nil
}
