package grpc

import (
	"log"
	"net"

	"github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
	//db "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/inmemory"
	db "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/postgre"
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
	var db s.Storage = db.DB {}
	calendarpb.RegisterCalendarEventsServer(grpcServer, calendar.GetEventServer(&db))
	err = grpcServer.Serve(lis)
	return err
}
