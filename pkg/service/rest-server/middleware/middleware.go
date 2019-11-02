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
	Name: "request_duration_milliseconds",
	Help: "The HTTP request latencies in milliseconds.",
}, []string{"method", "route", "status"})

var requestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "request_http_requests_total",
	Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
}, []string{"code", "method"})

//nolint
func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestTotal)
}

func Instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		duration := float64(m.Duration / time.Millisecond)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		requestDuration.WithLabelValues(
			r.Method, routePattern, strconv.Itoa(m.Code),
		).Observe(duration)
		requestTotal.WithLabelValues(strconv.Itoa(m.Code), r.Method).Add(1)
	})
}
