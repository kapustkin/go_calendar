package grpc

import (
	"log"
	"net"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run запуск GRPC сервера
func Run(addres string) error {
	lis, err := net.Listen("tcp", addres)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	calendarpb.RegisterCalendarEventsServer(grpcServer, calendar.GetEventServer())
	err = grpcServer.Serve(lis)
	return err
}
