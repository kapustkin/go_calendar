package grpc

import (
	"fmt"
	"net"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	s "github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/inmemory"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/postgre"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run запуск GRPC сервера
func Run() error {
	c := config.InitConfig()

	var db s.Storage
	switch c.StorageType {
	case 0:
		db = inmemory.DB{}
		db.Init(c.ConnectionString)
	case 1:
		db = postgre.DB{}
		db.Init(c.ConnectionString)
	default:
		return fmt.Errorf("storage type %d not supported", c.StorageType)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return fmt.Errorf("grpc failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	calendarpb.RegisterCalendarEventsServer(grpcServer, calendar.GetEventServer(&db))
	err = grpcServer.Serve(lis)
	return err
}
