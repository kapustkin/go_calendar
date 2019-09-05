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
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/logger"
	"github.com/sirupsen/logrus"
)

// Run основной обработчик
func Run(args []string) error {
	c := config.InitConfig()
	// Data Access Layer инициализация
	dal.Init(c.GRPC)

	r := chi.NewRouter()
	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Logging
	switch c.Logging {
	case 1:
		r.Use(middleware.Logger)
	case 2:
		log := logrus.New()
		log.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
		r.Use(logger.NewChiLogger(log))
	default:
		log.Printf("Warning! Starting without logging... \n")
	}

	// Routes
	r.Route("/calendar", func(r chi.Router) {
		r.Get("/{user}", calendar.GetEvents)
		r.Post("/{user}/add", calendar.AddEvent)
		r.Post("/{user}/edit", calendar.EditEvent)
		r.Post("/{user}/remove", calendar.RemoveEvent)
	})

	return http.ListenAndServe(fmt.Sprintf("%s:%v", c.Host, c.Port), r)
}
