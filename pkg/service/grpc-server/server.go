package grpc

import (
	"net"

	calendarpb "github.com/kapustkin/go_calendar/pkg/api/v1"
	"github.com/kapustkin/go_calendar/pkg/logger"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/handlers/calendar"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage"
	log "github.com/sirupsen/logrus"

	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/inmemory"
	"github.com/kapustkin/go_calendar/pkg/service/grpc-server/storage/postgre"
	"google.golang.org/grpc"
)

// Run запуск GRPC сервера
func Run() error {
	logger.Init("grpc-server", "0.0.1")
	log.Info("starting app...")
	conf := config.InitConfig()
	log.Infof("use config: %v", conf)
	db := getStorage(conf.StorageType, conf.ConnectionString)

	lis, err := net.Listen("tcp", conf.Host)
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
