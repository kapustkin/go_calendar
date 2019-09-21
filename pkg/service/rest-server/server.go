package rest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/dal"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/handlers/calendar"

	//"github.com/kapustkin/go_calendar/pkg/service/rest-server/logger"
	"github.com/kapustkin/go_calendar/pkg/logger"
)

// Run основной обработчик
func Run(args []string) error {
	// logger init
	applogger := logger.Init("rest-server", "0.0.1")
	applogger.Info("starting app...")
	c := config.InitConfig()
	applogger.Infof("use config: %v", c)
	// data access Layer init
	grpcDal := dal.Init(c.GRPC)

	r := chi.NewRouter()
	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	// Logging
	switch c.Logging {
	case 1:
		r.Use(middleware.Logger)
	case 2:
		r.Use(logger.NewChiLogger())
	default:
		log.Printf("Warning! Starting without logging... \n")
	}

	calendarService := calendar.Init(grpcDal)

	// Routes
	r.Route("/calendar", func(r chi.Router) {
		r.Get("/{user}", calendarService.GetEvents)
		r.Post("/{user}/add", calendarService.AddEvent)
		r.Post("/{user}/edit", calendarService.EditEvent)
		r.Post("/{user}/remove", calendarService.RemoveEvent)
	})
	applogger.Infof("listner started...")
	return http.ListenAndServe(fmt.Sprintf("%s:%v", c.Host, c.Port), r)
}
