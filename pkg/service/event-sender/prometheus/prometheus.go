package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//nolint
var regCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "event_sender_message_counter",
	Help: "кол-во отправленных оповещений",
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
