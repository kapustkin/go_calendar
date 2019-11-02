package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kapustkin/go_calendar/pkg/logger"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/config"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/dal"
	"github.com/kapustkin/go_calendar/pkg/service/rest-server/handlers/calendar"
	prometeus "github.com/kapustkin/go_calendar/pkg/service/rest-server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Run основной обработчик
func Run(args []string) error {
	// logger init
	logger.Init("rest-server", "0.0.1")
	log.Info("starting app...")
	conf := config.InitConfig()
	log.Infof("use config: %v", conf)
	// data access Layer init
	grpcDal := dal.Init(conf.GRPC)

	r := chi.NewRouter()
	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// prometeus middleware
	r.Use(prometeus.Instrument)
	r.Use(middleware.Timeout(60 * time.Second))
	// Logging
	switch conf.Logging {
	case 1:
		r.Use(middleware.Logger)
	case 2:
		r.Use(logger.NewChiLogger())
	default:
		log.Warn("starting without request logging...")
	}

	calendarService := calendar.Init(grpcDal)

	// Healthchecks
	r.Route("/", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("OK"))
			if err != nil {
				log.Fatal(err)
			}
		})
	})

	// Routes
	r.Route("/calendar", func(r chi.Router) {
		r.Get("/{user}", calendarService.GetEvents)
		r.Post("/{user}/add", calendarService.AddEvent)
		r.Post("/{user}/edit", calendarService.EditEvent)
		r.Post("/{user}/remove", calendarService.RemoveEvent)
	})
	log.Infof("listner started...")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Fatalf("metrics listener error %v", err)
		}
	}()

	return http.ListenAndServe(conf.Host, r)
}
