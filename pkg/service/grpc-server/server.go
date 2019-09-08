package grpc

import (
	"fmt"
	"net"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/inmemory"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/postgre"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run запуск GRPC сервера
func Run() error {
	conf := config.InitConfig()

	db, err := getStorage(conf.StorageType, conf.ConnectionString)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return fmt.Errorf("grpc failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	calendarpb.RegisterCalendarEventsServer(grpcServer, calendar.GetEventServer(db))
	err = grpcServer.Serve(lis)
	return err
}

func getStorage(storageType int, connectionString string) (*storage.Storage, error) {
	var db storage.Storage
	switch storageType {
	case 0:
		db = inmemory.DB{}
		db.Init(connectionString)
	case 1:
		db = postgre.DB{}
		db.Init(connectionString)
	default:
		return nil, fmt.Errorf("storage type %d not supported", storageType)
	}

	return &db, nil
}
