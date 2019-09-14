package grpc

import (
	"fmt"
	"log"
	"net"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/inmemory"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/postgre"
	"google.golang.org/grpc"
)

// Run запуск GRPC сервера
func Run() error {
	conf := config.InitConfig()

	db := getStorage(conf.StorageType, conf.ConnectionString)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	calendarpb.RegisterCalendarEventsServer(grpcServer, calendar.GetEventServer(db))
	err = grpcServer.Serve(lis)
	return err
}

func getStorage(storageType int, connectionString string) *storage.Storage {
	switch storageType {
	case 0:
		var db storage.Storage
		db = inmemory.Init()
		return &db
	case 1:
		var db storage.Storage
		db = postgre.Init(connectionString)
		return &db
	default:
		log.Panicf("storage type %d not supported", storageType)
	}
	return nil
}
