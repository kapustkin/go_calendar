package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//nolint
var regCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "calendar",
	Subsystem: "eventsender",
	Name:      "messages",
	Help:      "message per time",
})

// Init запуск мониторинга
func Init() {
	prometheus.MustRegister(regCounter)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Fatalf("metrics listener error %v", err)
		}
	}()
}

func RegisterEvent() {
	regCounter.Inc()
}
