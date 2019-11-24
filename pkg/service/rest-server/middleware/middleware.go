package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
)

var requestDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Namespace: "calendar",
	Subsystem: "restserver",
	Name:      "request_duration_milliseconds",
	Help:      "The HTTP request latencies in milliseconds.",
}, []string{"method", "route", "status"})

var requestCounter = prometheus.NewHistogram(prometheus.HistogramOpts{
	Namespace: "calendar",
	Subsystem: "restserver",
	Name:      "requests",
	Help:      "Operation per time",
	Buckets:   prometheus.LinearBuckets(0, 10, 20),
})

//nolint
func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestCounter)
}

func Monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		duration := float64(m.Duration / time.Millisecond)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		requestDuration.WithLabelValues(
			r.Method, routePattern, strconv.Itoa(m.Code),
		).Observe(duration)

		requestCounter.Observe(1)
	})
}
